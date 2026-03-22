import { useEffect, useMemo, useState } from 'react';
import { GitBranch, Rocket, X } from 'lucide-react';
import { ArmsHttpError } from '../../api/armsClient';
import type { ApiProductDetail, ApiVersion } from '../../api/armsTypes';
import type { Task } from '../../domain/types';
import { formatDurationMs } from '../../lib/time';

export type MissionControlWorkspaceStats = {
  thisWeek: number;
  inProgress: number;
  total: number;
  completionPct: number;
};

type Props = {
  open: boolean;
  onClose: () => void;
  workspaceName: string;
  workspaceIcon: string;
  isOnline: boolean;
  fetchVersion: () => Promise<ApiVersion>;
  productDetail: ApiProductDetail | null;
  productTasks: Task[];
  workspaceStats: MissionControlWorkspaceStats;
};

const SESSION_KEY = 'fishtank-session-start-ms';

function sessionStartMs(): number {
  try {
    const raw = sessionStorage.getItem(SESSION_KEY);
    if (raw) {
      const n = parseInt(raw, 10);
      if (Number.isFinite(n)) return n;
    }
    const now = Date.now();
    sessionStorage.setItem(SESSION_KEY, String(now));
    return now;
  } catch {
    return Date.now();
  }
}

function shortSha(commit: string | undefined): string {
  const c = commit?.trim() ?? '';
  if (c.length <= 12) return c || '—';
  return `${c.slice(0, 7)}…`;
}

function collectTaskPaths(tasks: Task[]): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const t of tasks) {
    for (const p of [t.worktreePath, t.sandboxPath]) {
      const v = p?.trim();
      if (!v || seen.has(v)) continue;
      seen.add(v);
      out.push(v);
      if (out.length >= 12) return out;
    }
  }
  return out;
}

export function MissionControlOverviewModal({
  open,
  onClose,
  workspaceName,
  workspaceIcon,
  isOnline,
  fetchVersion,
  productDetail,
  productTasks,
  workspaceStats,
}: Props) {
  const [tick, setTick] = useState(() => Date.now());
  const [versionLoading, setVersionLoading] = useState(false);
  const [versionError, setVersionError] = useState<string | null>(null);
  const [version, setVersion] = useState<ApiVersion | null>(null);

  const sessionStart = useMemo(() => sessionStartMs(), []);

  useEffect(() => {
    if (!open) return;
    const id = window.setInterval(() => setTick(Date.now()), 1000);
    return () => window.clearInterval(id);
  }, [open]);

  useEffect(() => {
    if (!open) return;
    setVersionError(null);
    setVersion(null);
    setVersionLoading(true);
    let cancelled = false;
    void (async () => {
      try {
        const v = await fetchVersion();
        if (!cancelled) setVersion(v);
      } catch (e) {
        if (!cancelled) {
          setVersionError(e instanceof ArmsHttpError ? e.message : e instanceof Error ? e.message : 'Could not load version');
        }
      } finally {
        if (!cancelled) setVersionLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [open, fetchVersion]);

  const taskPaths = useMemo(() => collectTaskPaths(productTasks), [productTasks]);

  if (!open) return null;

  const uptimeLabel = formatDurationMs(tick - sessionStart);
  const s = workspaceStats;

  return (
    <div className="ft-modal-root" role="dialog" aria-modal="true" aria-labelledby="ft-mc-overview-title">
      <button type="button" className="ft-modal-backdrop" aria-label="Close" onClick={onClose} />
      <div className="ft-modal-panel ft-mc-overview-modal">
        <div className="ft-modal-head">
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', minWidth: 0 }}>
            <span className="ft-mc-overview-modal-icon" aria-hidden>
              <Rocket size={20} />
            </span>
            <h2 id="ft-mc-overview-title" className="ft-mc-overview-modal-title">
              Mission Control
            </h2>
          </div>
          <button type="button" className="ft-btn-icon" title="Close" aria-label="Close" onClick={onClose}>
            <X size={18} />
          </button>
        </div>

        <div className="ft-modal-body ft-mc-overview-body">
          <section className="ft-mc-overview-section">
            <h3 className="ft-mc-overview-label">Workspace</h3>
            <div className="ft-mc-overview-workspace">
              <span aria-hidden className="ft-mc-overview-ws-icon">
                {workspaceIcon}
              </span>
              <span className="ft-truncate" style={{ fontWeight: 600 }}>
                {workspaceName}
              </span>
            </div>
          </section>

          <section className="ft-mc-overview-section">
            <h3 className="ft-mc-overview-label">Connection</h3>
            <dl className="ft-mc-overview-dl">
              <div>
                <dt>Online status</dt>
                <dd>
                  <span className={isOnline ? 'ft-mc-overview-pill ft-mc-overview-pill--on' : 'ft-mc-overview-pill ft-mc-overview-pill--off'}>
                    {isOnline ? 'Online' : 'Offline'}
                  </span>
                </dd>
              </div>
              <div>
                <dt>Session uptime</dt>
                <dd className="ft-mono">{uptimeLabel}</dd>
              </div>
            </dl>
            <p className="ft-muted" style={{ fontSize: '0.7rem', marginTop: '0.35rem', lineHeight: 1.4 }}>
              Uptime is measured in this browser tab since the session started (not arms server uptime).
            </p>
          </section>

          <section className="ft-mc-overview-section" aria-labelledby="ft-mc-overview-stats">
            <h3 id="ft-mc-overview-stats" className="ft-mc-overview-label">
              Task stats (this product)
            </h3>
            <div className="ft-mc-overview-stats" role="group">
              <div className="ft-mc-overview-stat">
                <span className="ft-mc-overview-stat-value ft-mc-overview-stat-value--green">{s.thisWeek}</span>
                <span className="ft-mc-overview-stat-label">This week</span>
              </div>
              <div className="ft-mc-overview-stat">
                <span className="ft-mc-overview-stat-value ft-mc-overview-stat-value--blue">{s.inProgress}</span>
                <span className="ft-mc-overview-stat-label">In progress</span>
              </div>
              <div className="ft-mc-overview-stat">
                <span className="ft-mc-overview-stat-value">{s.total}</span>
                <span className="ft-mc-overview-stat-label">Total</span>
              </div>
              <div className="ft-mc-overview-stat">
                <span className="ft-mc-overview-stat-value">{s.completionPct}%</span>
                <span className="ft-mc-overview-stat-label">Done</span>
                <span className="ft-mc-overview-stat-bar" aria-hidden>
                  <span className="ft-mc-overview-stat-bar-fill" style={{ width: `${s.completionPct}%` }} />
                </span>
              </div>
            </div>
          </section>

          <section className="ft-mc-overview-section">
            <h3 className="ft-mc-overview-label">Repositories &amp; git</h3>
            <p className="ft-muted" style={{ fontSize: '0.72rem', marginBottom: '0.65rem', lineHeight: 1.45 }}>
              arms does not expose live <code className="ft-mono">git status</code> for clones; below is the API server build and
              configured product repo. Task paths are local hints from loaded tasks.
            </p>

            <div className="ft-mc-overview-repo-card">
              <div className="ft-mc-overview-repo-head">
                <GitBranch size={14} aria-hidden />
                <span>arms API (server build)</span>
              </div>
              {versionLoading ? (
                <p className="ft-muted" style={{ fontSize: '0.8rem' }}>
                  Loading…
                </p>
              ) : versionError ? (
                <p className="ft-banner ft-banner--error" style={{ fontSize: '0.75rem', padding: '0.35rem 0.5rem' }}>
                  {versionError}
                </p>
              ) : version ? (
                <ul className="ft-mc-overview-repo-meta">
                  <li>
                    <span className="ft-muted">Commit</span> <code className="ft-mono">{shortSha(version.commit)}</code>
                    {version.dirty ? (
                      <span className="ft-mc-overview-dirty" title="Working tree dirty at build">
                        dirty
                      </span>
                    ) : null}
                  </li>
                  <li>
                    <span className="ft-muted">Tag / version</span> {version.tag || version.number || version.version || '—'}
                  </li>
                  {version.commits_after_tag != null ? (
                    <li>
                      <span className="ft-muted">Commits after tag</span> {version.commits_after_tag}
                    </li>
                  ) : null}
                </ul>
              ) : (
                <p className="ft-muted">—</p>
              )}
            </div>

            <div className="ft-mc-overview-repo-card">
              <div className="ft-mc-overview-repo-head">
                <GitBranch size={14} aria-hidden />
                <span>Product repository (configured)</span>
              </div>
              {productDetail?.repo_url || productDetail?.repo_branch ? (
                <ul className="ft-mc-overview-repo-meta">
                  {productDetail.repo_url ? (
                    <li>
                      <span className="ft-muted">Remote</span>{' '}
                      <a href={productDetail.repo_url} className="ft-mc-overview-link" target="_blank" rel="noreferrer">
                        {productDetail.repo_url}
                      </a>
                    </li>
                  ) : null}
                  {productDetail.repo_branch ? (
                    <li>
                      <span className="ft-muted">Branch</span> <code className="ft-mono">{productDetail.repo_branch}</code>
                    </li>
                  ) : null}
                </ul>
              ) : (
                <p className="ft-muted" style={{ fontSize: '0.8rem' }}>
                  No repo URL or branch on this product record.
                </p>
              )}
            </div>

            {taskPaths.length > 0 ? (
              <div className="ft-mc-overview-repo-card">
                <div className="ft-mc-overview-repo-head">
                  <GitBranch size={14} aria-hidden />
                  <span>Task-linked paths</span>
                </div>
                <ul className="ft-mc-overview-path-list">
                  {taskPaths.map((p) => (
                    <li key={p}>
                      <code className="ft-mono">{p}</code>
                    </li>
                  ))}
                </ul>
              </div>
            ) : null}
          </section>
        </div>
      </div>
    </div>
  );
}
