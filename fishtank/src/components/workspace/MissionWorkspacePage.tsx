import { useMemo, useState } from 'react';
import { LayoutGrid, Radio, Users } from 'lucide-react';
import { useMissionUi } from '../../context/MissionUiContext';
import { WorkspaceHeaderBar } from '../shell/WorkspaceHeaderBar';
import { AgentsPanel } from './AgentsPanel';
import { LiveFeedPanel } from './LiveFeedPanel';
import { MissionControlSidebar } from './MissionControlSidebar';
import { MissionQueuePanel } from './MissionQueuePanel';

type MobileTab = 'queue' | 'agents' | 'feed';

const WEEK_MS = 7 * 24 * 60 * 60 * 1000;

/** Desktop: Mission Control shell (nav + board + live activity); mobile uses tabs. */
export function MissionWorkspacePage() {
  const { apiError, dismissError, activeWorkspace, tasks } = useMissionUi();
  const [tab, setTab] = useState<MobileTab>('queue');
  const [boardSearch, setBoardSearch] = useState('');
  const [assigneeAgentId, setAssigneeAgentId] = useState<string | null>(null);
  const [newTaskOpen, setNewTaskOpen] = useState(false);
  const [agentsPaused, setAgentsPaused] = useState(false);
  const [leftPanel, setLeftPanel] = useState<'nav' | 'agents'>('nav');

  const stats = useMemo(() => {
    if (!activeWorkspace) {
      return { thisWeek: 0, inProgress: 0, total: 0, completionPct: 0 };
    }
    const scoped = tasks.filter((t) => t.workspaceId === activeWorkspace.id);
    const now = Date.now();
    const thisWeek = scoped.filter((t) => now - new Date(t.updatedAt).getTime() < WEEK_MS).length;
    const inProgress = scoped.filter((t) =>
      t.status === 'in_progress' || t.status === 'testing' || t.status === 'convoy_active',
    ).length;
    const total = scoped.length;
    const done = scoped.filter((t) => t.status === 'done').length;
    const completionPct = total === 0 ? 0 : Math.round((done / total) * 100);
    return { thisWeek, inProgress, total, completionPct };
  }, [tasks, activeWorkspace]);

  const queueProps = {
    boardSearch,
    assigneeAgentId,
    onAssigneeAgentIdChange: setAssigneeAgentId,
    newTaskOpen,
    onNewTaskOpenChange: setNewTaskOpen,
  };

  return (
    <div className="ft-screen-fixed">
      <WorkspaceHeaderBar
        missionControl={{
          boardSearch,
          onBoardSearchChange: setBoardSearch,
          agentsPaused,
          onAgentsPausedToggle: () => setAgentsPaused((p) => !p),
          workspaceStats: stats,
        }}
      />
      {apiError ? (
        <div className="ft-container" style={{ padding: '0.35rem 1rem 0' }}>
          <div
            className="ft-banner ft-banner--error"
            role="alert"
            style={{ display: 'flex', justifyContent: 'space-between', gap: '0.5rem', alignItems: 'center' }}
          >
            <span style={{ fontSize: '0.8rem' }}>{apiError}</span>
            <button type="button" className="ft-btn-ghost" style={{ fontSize: '0.75rem' }} onClick={dismissError}>
              Dismiss
            </button>
          </div>
        </div>
      ) : null}

      <div className="ft-desktop-only ft-mc-desktop-row">
        <MissionControlSidebar stats={stats} leftPanel={leftPanel} onLeftPanelChange={setLeftPanel} />
        <MissionQueuePanel {...queueProps} />
        <LiveFeedPanel variant="activity" />
      </div>

      <div className="ft-mobile-only" style={{ flex: 1, minHeight: 0 }}>
        <div className="ft-mc-stats-row ft-mc-stats-row--mobile" role="group" aria-label="Workspace stats">
          <div className="ft-mc-stat-pill ft-mc-stat-pill--green">
            <span className="ft-mc-stat-pill-value">{stats.thisWeek}</span>
            <span className="ft-mc-stat-pill-label">Week</span>
          </div>
          <div className="ft-mc-stat-pill ft-mc-stat-pill--blue">
            <span className="ft-mc-stat-pill-value">{stats.inProgress}</span>
            <span className="ft-mc-stat-pill-label">Active</span>
          </div>
          <div className="ft-mc-stat-pill">
            <span className="ft-mc-stat-pill-value">{stats.total}</span>
            <span className="ft-mc-stat-pill-label">Total</span>
          </div>
          <div className="ft-mc-stat-pill ft-mc-stat-pill--progress">
            <span className="ft-mc-stat-pill-value">{stats.completionPct}%</span>
            <span className="ft-mc-stat-pill-label">Done</span>
          </div>
        </div>
        <div className="ft-mobile-tab-bar" role="tablist" aria-label="Workspace panels">
          <button
            type="button"
            role="tab"
            aria-selected={tab === 'queue'}
            className={`ft-mobile-tab ${tab === 'queue' ? 'ft-mobile-tab--active' : ''}`}
            onClick={() => setTab('queue')}
          >
            <LayoutGrid size={14} style={{ display: 'block', margin: '0 auto 0.2rem' }} aria-hidden />
            Tasks
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={tab === 'agents'}
            className={`ft-mobile-tab ${tab === 'agents' ? 'ft-mobile-tab--active' : ''}`}
            onClick={() => setTab('agents')}
          >
            <Users size={14} style={{ display: 'block', margin: '0 auto 0.2rem' }} aria-hidden />
            Agents
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={tab === 'feed'}
            className={`ft-mobile-tab ${tab === 'feed' ? 'ft-mobile-tab--active' : ''}`}
            onClick={() => setTab('feed')}
          >
            <Radio size={14} style={{ display: 'block', margin: '0 auto 0.2rem' }} aria-hidden />
            Activity
          </button>
        </div>
        <div style={{ flex: 1, minHeight: 0, overflow: 'hidden', display: 'flex', flexDirection: 'column' }} role="tabpanel">
          {tab === 'queue' ? <MissionQueuePanel {...queueProps} /> : null}
          {tab === 'agents' ? <AgentsPanel /> : null}
          {tab === 'feed' ? <LiveFeedPanel variant="activity" /> : null}
        </div>
      </div>
    </div>
  );
}
