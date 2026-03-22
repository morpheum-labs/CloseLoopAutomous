import { useEffect } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import { MarkdownReadViewer } from './markdownRead/MarkdownReadViewer';
import 'katex/dist/katex.min.css';
import 'highlight.js/styles/github-dark.css';

export type MarkdownReadModalProps = {
  open: boolean;
  onClose: () => void;
  title: string;
  content: string;
};

export function MarkdownReadModal({ open, onClose, title, content }: MarkdownReadModalProps) {
  useEffect(() => {
    if (!open) return undefined;
    const prev = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = prev;
    };
  }, [open]);

  useEffect(() => {
    if (!open) return undefined;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [open, onClose]);

  if (!open) return null;

  return createPortal(
    <div className="ft-docs-read-overlay" role="dialog" aria-modal="true" aria-label={title}>
      <button
        type="button"
        className="ft-docs-read-backdrop"
        aria-label="Close reader"
        onClick={onClose}
      />
      <div className="ft-docs-read-panel">
        <header className="ft-docs-read-header">
          <span className="ft-docs-read-title">{title}</span>
          <button
            type="button"
            className="ft-btn-ghost ft-docs-read-close"
            onClick={onClose}
            aria-label="Close"
          >
            <X size={20} aria-hidden />
          </button>
        </header>
        <div className="ft-docs-read-body">
          <MarkdownReadViewer content={content} className="ft-docs-read-viewer" />
        </div>
      </div>
    </div>,
    document.body,
  );
}
