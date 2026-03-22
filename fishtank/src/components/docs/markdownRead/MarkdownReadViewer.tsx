/**
 * Markdown rendering (marked + KaTeX + highlight.js). Mermaid code fences are shown as source only:
 * bundling mermaid pulls khroma, which breaks under Bun (`adjust is not a function`). Use markdownApp for diagrams.
 */
import React, { useEffect, useCallback } from 'react';
import { cn, convertLatexDelimiters, convertLatexToKatex } from './md-utils';
import { Marked, Renderer, type Tokens } from 'marked';
import hljs from 'highlight.js';
import { markedHighlight } from 'marked-highlight';
import markedKatex from 'marked-katex-extension';
import markedCodePreview from 'marked-code-preview';

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

const katexOptions = {
  throwOnError: false,
  output: 'html' as const,
  displayMode: false,
  minRuleThickness: 0.15,
  maxExpand: 900,
  maxSize: 100000,
  nonStandard: true,
};

const createMarkedInstance = (customRenderer: Renderer) => {
  const instance = new Marked(
    markedHighlight({
      emptyLangClass: 'hljs',
      langPrefix: 'hljs language-',
      highlight(code, lang) {
        if (lang === 'mermaid') {
          return escapeHtml(code);
        }
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
      },
    }),
    {
      renderer: customRenderer,
      breaks: true,
      gfm: true,
      pedantic: false,
    },
  );

  instance.use(markedKatex(katexOptions));
  instance.use(markedCodePreview());
  return instance;
};

const createCustomRenderer = (): Renderer => {
  const renderer = new Renderer();

  const inlineMarked = new Marked({
    breaks: true,
    gfm: true,
    pedantic: false,
    renderer,
  });
  inlineMarked.use(markedKatex(katexOptions));
  inlineMarked.use(markedCodePreview());

  const tempMarked = new Marked(
    markedHighlight({
      emptyLangClass: 'hljs',
      langPrefix: 'hljs language-',
      highlight(code, lang) {
        if (lang === 'mermaid') {
          return escapeHtml(code);
        }
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
      },
    }),
    {
      breaks: true,
      gfm: true,
      pedantic: false,
      renderer,
    },
  );
  tempMarked.use(markedKatex(katexOptions));
  tempMarked.use(markedCodePreview());

  renderer.paragraph = function (text: Tokens.Paragraph): string {
    const textContent = String(text.raw || text);
    const result = inlineMarked.parseInline(textContent);
    return `<p class="story">${result}</p>`;
  };

  renderer.text = function (text: Tokens.Text | Tokens.Escape) {
    const textContent = String(text.text || text);
    return textContent;
  };

  renderer.list = function (list: Tokens.List): string {
    const type = list.ordered ? 'ol' : 'ul';
    const startAttr = list.ordered && list.start !== 1 ? ` start="${list.start}"` : '';
    return `<${type}${startAttr} class="L1">${list.items.map((item) => this.listitem(item)).join('')}</${type}>`;
  };

  renderer.listitem = function (item: Tokens.ListItem): string {
    let processedContent = '';
    if (item.text) {
      const sample_context = item.text;
      try {
        const result = tempMarked.parseInline(sample_context);
        processedContent = typeof result === 'string' ? result : sample_context;
      } catch {
        processedContent = sample_context;
      }
    } else if (item.tokens && item.tokens.length > 0) {
      processedContent = item.tokens
        .map((token) => {
          const rendererMethod = (this as unknown as Record<string, (t: typeof token) => string>)[token.type];
          return rendererMethod ? rendererMethod.call(this, token) : token.raw || '';
        })
        .join('');
    }

    return `<li class="L1i">${processedContent}</li>\n`;
  };

  renderer.strong = function (text: Tokens.Strong): string {
    const textContent = String(text.text || text);
    return `<strong>${textContent}</strong>`;
  };

  renderer.code = function (code: Tokens.Code): string {
    if (code.lang === 'mermaid') {
      const codeString = typeof code === 'string' ? code : code.text || String(code);
      return `<pre class="ft-docs-md-read__mermaid-src"><code class="hljs language-mermaid">${escapeHtml(codeString)}</code></pre>`;
    }
    return this.constructor.prototype.code.call(this, code);
  };

  renderer.tablecell = function (cell: Tokens.TableCell): string {
    const tag = cell.header ? 'th' : 'td';
    const sample_context = cell.text;
    const result = tempMarked.parseInline(sample_context);
    return `<${tag} class="L2c">${result}</${tag}>`;
  };

  return renderer;
};

interface MarkdownReadViewerProps {
  content: string;
  className?: string;
}

export function MarkdownReadViewer({ content, className }: MarkdownReadViewerProps) {
  const [processedHtml, setProcessedHtml] = React.useState('');
  const [isMounted, setIsMounted] = React.useState(false);

  const customRenderer = React.useMemo(() => createCustomRenderer(), []);

  const processMarkdown = useCallback((contentToProcess: string, renderer: Renderer) => {
    try {
      if (!contentToProcess.trim()) {
        return '';
      }
      let processedContent = contentToProcess;
      processedContent = processedContent.replace(/^(\s+)(\\[\[\(][\s\S]*?\\[\]\)])/gm, (match, indent: string, math: string) => {
        const effectiveIndent = indent.replace(/\t/g, '    ').length;
        if (effectiveIndent >= 4) {
          return math;
        }
        return match;
      });

      processedContent = convertLatexDelimiters(processedContent);
      processedContent = convertLatexToKatex(processedContent);

      const currencyPlaceholders: string[] = [];
      const currencyRegex = /\$(\d+(?:\.\d+)?[MKB])(?=\s|$|[^\w])/g;
      processedContent = processedContent.replace(currencyRegex, (match) => {
        const placeholder = `CURRENCY_PLACEHOLDER_${currencyPlaceholders.length}`;
        currencyPlaceholders.push(match);
        return placeholder;
      });

      const markedInstance = createMarkedInstance(renderer);
      let renderedHtml = markedInstance.parse(processedContent) as string;

      currencyPlaceholders.forEach((currency, index) => {
        const placeholder = `CURRENCY_PLACEHOLDER_${index}`;
        renderedHtml = renderedHtml.replace(new RegExp(placeholder, 'g'), currency);
      });

      return renderedHtml;
    } catch (error) {
      console.error('Error processing markdown:', error);
      return `<div class="error">${error instanceof Error ? escapeHtml(error.message) : 'Unknown error'}</div>`;
    }
  }, []);

  useEffect(() => {
    setIsMounted(true);
    return () => {
      setIsMounted(false);
    };
  }, []);

  useEffect(() => {
    if (isMounted && content && content.trim()) {
      setProcessedHtml(processMarkdown(content, customRenderer));
    } else if (!content || !content.trim()) {
      setProcessedHtml('');
    }
  }, [content, isMounted, customRenderer, processMarkdown]);

  return (
    <div className={cn('ft-docs-md-read__card overflow-auto', className)}>
      <div className="markdown-body" dangerouslySetInnerHTML={{ __html: processedHtml }} />
    </div>
  );
}
