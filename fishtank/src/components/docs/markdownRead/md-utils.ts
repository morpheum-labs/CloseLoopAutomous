import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}

/** Convert LaTeX delimiters \( \) and \[ \] to $ / $$ for KaTeX extension. */
export function convertLatexDelimiters(text: string): string {
  const inlineMathRegex = /\\\([\s\S]*?\\\)/g;
  let result = text.replace(inlineMathRegex, (match) => {
    const content = match.slice(2, -2);
    return `$${content}$`;
  });
  const blockMathRegex = /\\\[[\s\S]*?\\\]/g;
  result = result.replace(blockMathRegex, (match) => {
    const content = match.slice(2, -2);
    return `$$${content}$$`;
  });
  return result;
}

/** Best-effort LaTeX → KaTeX-friendly substitutions (from markdownApp). */
export function convertLatexToKatex(text: string): string {
  let result = text;
  const latexToKatexConversions: { from: RegExp; to: string }[] = [
    { from: /\\mbox\{([^}]+)\}/g, to: '\\text{$1}' },
    { from: /\\textrm\{([^}]+)\}/g, to: '\\text{$1}' },
    { from: /\\dfrac\{([^}]+)\}\{([^}]+)\}/g, to: '\\frac{$1}{$2}' },
    { from: /\\tfrac\{([^}]+)\}\{([^}]+)\}/g, to: '\\frac{$1}{$2}' },
    { from: /\\cfrac\{([^}]+)\}\{([^}]+)\}/g, to: '\\frac{$1}{$2}' },
    { from: /\\usepackage\{[^}]+\}/g, to: '' },
    { from: /\\begin\{document\}/g, to: '' },
    { from: /\\end\{document\}/g, to: '' },
  ];
  for (const { from, to } of latexToKatexConversions) {
    result = result.replace(from, to);
  }
  result = result.replace(/\\text\{([^}]+)\}/g, (match, content: string) => {
    return `\\text{${content.replace(/_/g, '\\_').replace(/%/g, '\\%')}}`;
  });
  return result;
}
