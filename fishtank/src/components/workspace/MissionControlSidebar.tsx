import { useNavigate } from 'react-router-dom';
import {
  Activity,
  Bot,
  BookOpen,
  Calendar,
  CheckSquare,
  ClipboardCheck,
  Factory,
  FileText,
  LayoutGrid,
  MessageSquare,
  Network,
  Radar,
  Rocket,
  Settings,
  Users,
  type LucideIcon,
} from 'lucide-react';
import { AgentsPanel } from './AgentsPanel';

type NavEntry =
  | { id: string; label: string; icon: LucideIcon; behavior: 'tasks' }
  | { id: string; label: string; icon: LucideIcon; behavior: 'agents' }
  | { id: string; label: string; icon: LucideIcon; behavior: 'route'; to: string }
  | { id: string; label: string; icon: LucideIcon; behavior: 'disabled' };

const NAV_ENTRIES: NavEntry[] = [
  { id: 'tasks', label: 'Tasks', icon: CheckSquare, behavior: 'tasks' },
  { id: 'agents', label: 'Agents', icon: Bot, behavior: 'agents' },
  { id: 'activity_log', label: 'Activity log', icon: Activity, behavior: 'route', to: '/activity' },
  { id: 'autopilot', label: 'Autopilot', icon: Rocket, behavior: 'route', to: '/autopilot' },
  { id: 'content', label: 'Content', icon: FileText, behavior: 'disabled' },
  { id: 'approvals', label: 'Approvals', icon: ClipboardCheck, behavior: 'disabled' },
  { id: 'council', label: 'Council', icon: Users, behavior: 'disabled' },
  { id: 'calendar', label: 'Calendar', icon: Calendar, behavior: 'disabled' },
  { id: 'projects', label: 'Projects', icon: LayoutGrid, behavior: 'disabled' },
  { id: 'memory', label: 'Memory', icon: Network, behavior: 'disabled' },
  { id: 'docs', label: 'Docs', icon: BookOpen, behavior: 'disabled' },
  { id: 'people', label: 'People', icon: Users, behavior: 'disabled' },
  { id: 'office', label: 'Office', icon: LayoutGrid, behavior: 'disabled' },
  { id: 'team', label: 'Team', icon: Users, behavior: 'disabled' },
  { id: 'system', label: 'System', icon: Settings, behavior: 'disabled' },
  { id: 'radar', label: 'Radar', icon: Radar, behavior: 'disabled' },
  { id: 'factory', label: 'Factory', icon: Factory, behavior: 'disabled' },
  { id: 'pipeline', label: 'Pipeline', icon: Network, behavior: 'disabled' },
  { id: 'feedback', label: 'Feedback', icon: MessageSquare, behavior: 'disabled' },
];

type Stats = {
  thisWeek: number;
  inProgress: number;
  total: number;
  completionPct: number;
};

type Props = {
  stats: Stats;
  leftPanel: 'nav' | 'agents';
  onLeftPanelChange: (p: 'nav' | 'agents') => void;
};

export function MissionControlSidebar({ stats, leftPanel, onLeftPanelChange }: Props) {
  const navigate = useNavigate();

  return (
    <aside className="ft-mc-sidebar" aria-label="Mission navigation">
      {leftPanel === 'agents' ? (
        <div className="ft-mc-sidebar-agents">
          <button type="button" className="ft-mc-nav-back" onClick={() => onLeftPanelChange('nav')}>
            ← Tasks
          </button>
          <AgentsPanel embedded />
        </div>
      ) : (
        <>
          <div className="ft-mc-stats-row" role="group" aria-label="Workspace stats">
            <div className="ft-mc-stat-pill ft-mc-stat-pill--green">
              <span className="ft-mc-stat-pill-value">{stats.thisWeek}</span>
              <span className="ft-mc-stat-pill-label">This week</span>
            </div>
            <div className="ft-mc-stat-pill ft-mc-stat-pill--blue">
              <span className="ft-mc-stat-pill-value">{stats.inProgress}</span>
              <span className="ft-mc-stat-pill-label">In progress</span>
            </div>
            <div className="ft-mc-stat-pill">
              <span className="ft-mc-stat-pill-value">{stats.total}</span>
              <span className="ft-mc-stat-pill-label">Total</span>
            </div>
            <div className="ft-mc-stat-pill ft-mc-stat-pill--progress" title="Share of tasks in Done">
              <span className="ft-mc-stat-pill-value">{stats.completionPct}%</span>
              <span className="ft-mc-stat-pill-label">Done</span>
              <span className="ft-mc-stat-pill-bar" aria-hidden>
                <span className="ft-mc-stat-pill-bar-fill" style={{ width: `${stats.completionPct}%` }} />
              </span>
            </div>
          </div>

          <nav className="ft-mc-nav" aria-label="Primary">
            <ul className="ft-mc-nav-list">
              {NAV_ENTRIES.map((entry) => {
                const Icon = entry.icon;
                const isActive =
                  (entry.behavior === 'tasks' && leftPanel === 'nav') || (entry.behavior === 'agents' && leftPanel === 'agents');

                function handleClick() {
                  if (entry.behavior === 'tasks') onLeftPanelChange('nav');
                  else if (entry.behavior === 'agents') onLeftPanelChange('agents');
                  else if (entry.behavior === 'route') navigate(entry.to);
                }

                const disabled = entry.behavior === 'disabled';

                return (
                  <li key={entry.id}>
                    <button
                      type="button"
                      className={`ft-mc-nav-item ${isActive ? 'ft-mc-nav-item--active' : ''}`}
                      onClick={handleClick}
                      disabled={disabled}
                      title={disabled ? 'Coming soon' : undefined}
                    >
                      <Icon size={16} aria-hidden className="ft-mc-nav-icon" />
                      <span>{entry.label}</span>
                    </button>
                  </li>
                );
              })}
            </ul>
          </nav>
        </>
      )}
    </aside>
  );
}
