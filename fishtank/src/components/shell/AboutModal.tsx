import { useEffect, useState } from 'react';
import { X } from 'lucide-react';
import type { ApiVersion } from '../../api/armsTypes';
import { ArmsHttpError } from '../../api/armsClient';

type Props = {
  open: boolean;
  onClose: () => void;
  fetchVersion: () => Promise<ApiVersion>;
};

function displayVersion(v: ApiVersion): string {
  const n = v.number?.trim();
  if (n) return n;
  const t = v.tag?.trim();
  if (t) return t;
  return v.version?.trim() || '—';
}

export function AboutModal({ open, onClose, fetchVersion }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [info, setInfo] = useState<ApiVersion | null>(null);

  useEffect(() => {
    if (!open) return;
    setError(null);
    setInfo(null);
    setLoading(true);
    let cancelled = false;
    void (async () => {
      try {
        const data = await fetchVersion();
        if (!cancelled) setInfo(data);
      } catch (e) {
        if (cancelled) return;
        if (e instanceof ArmsHttpError) {
          setError(e.message);
        } else {
          setError(e instanceof Error ? e.message : 'Could not load version');
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [open, fetchVersion]);

  if (!open) return null;

  return (
    <div className="ft-modal-root" role="dialog" aria-modal="true" aria-labelledby="ft-about-title">
      <button type="button" className="ft-modal-backdrop" aria-label="Close" onClick={onClose} />
      <div className="ft-modal-panel" style={{ width: 'min(100%, 440px)' }}>
        <div className="ft-modal-head">
          <h2 id="ft-about-title" style={{ margin: 0, fontSize: '1.1rem', fontWeight: 600 }}>
            About Fishtank
          </h2>
          <button type="button" className="ft-btn-icon" onClick={onClose} aria-label="Close dialog">
            <X size={18} />
          </button>
        </div>
        <div className="ft-modal-body">
          <p className="ft-muted" style={{ margin: 0, fontSize: '0.875rem', lineHeight: 1.5 }}>
            Mission Control UI for arms. Backend build metadata from{' '}
            <code className="ft-mono">GET /api/version</code>.
          </p>

          {loading ? <p className="ft-muted" style={{ margin: 0 }}>Loading version…</p> : null}
          {error ? (
            <p className="ft-banner ft-banner--error" role="alert" style={{ margin: 0 }}>
              {error}
            </p>
          ) : null}

          {info && !loading ? (
            <dl
              style={{
                margin: 0,
                display: 'grid',
                gap: '0.65rem',
                fontSize: '0.875rem',
              }}
            >
              <div>
                <dt className="ft-field-label" style={{ marginBottom: '0.2rem' }}>
                  Arms version
                </dt>
                <dd style={{ margin: 0, fontWeight: 700, fontSize: '1.35rem', letterSpacing: '-0.02em' }}>
                  {displayVersion(info)}
                </dd>
              </div>
              <div>
                <dt className="ft-field-label" style={{ marginBottom: '0.2rem' }}>
                  Describe
                </dt>
                <dd style={{ margin: 0, wordBreak: 'break-all' }} className="ft-mono">
                  {info.version || '—'}
                </dd>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
                <div>
                  <dt className="ft-field-label" style={{ marginBottom: '0.2rem' }}>
                    Tag
                  </dt>
                  <dd style={{ margin: 0 }} className="ft-mono">
                    {info.tag || '—'}
                  </dd>
                </div>
                <div>
                  <dt className="ft-field-label" style={{ marginBottom: '0.2rem' }}>
                    Commit
                  </dt>
                  <dd style={{ margin: 0 }} className="ft-mono">
                    {info.commit || '—'}
                  </dd>
                </div>
              </div>
              {(info.commits_after_tag > 0 || info.dirty) && (
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', alignItems: 'center' }}>
                  {info.commits_after_tag > 0 ? (
                    <span className="ft-chip" style={{ fontSize: '0.75rem' }}>
                      +{info.commits_after_tag} commit{info.commits_after_tag === 1 ? '' : 's'} after tag
                    </span>
                  ) : null}
                  {info.dirty ? (
                    <span className="ft-chip" style={{ fontSize: '0.75rem' }}>
                      dirty working tree
                    </span>
                  ) : null}
                </div>
              )}
            </dl>
          ) : null}

          <div className="ft-modal-actions">
            <button type="button" className="ft-btn-primary" onClick={onClose}>
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
