import { GripVertical } from 'lucide-react';
import type { Agent, Task, TaskStatus } from '../../domain/types';
import { formatRelativeTime } from '../../lib/time';

type Props = {
  task: Task;
  agents: Agent[];
  onOpen: () => void;
};

function statusDotClass(status: TaskStatus): string {
  switch (status) {
    case 'done':
      return 'ft-task-dot ft-task-dot--green';
    case 'failed':
      return 'ft-task-dot ft-task-dot--red';
    case 'review':
      return 'ft-task-dot ft-task-dot--orange';
    case 'in_progress':
    case 'testing':
    case 'convoy_active':
      return 'ft-task-dot ft-task-dot--blue';
    default:
      return 'ft-task-dot ft-task-dot--muted';
  }
}

function taskTag(task: Task): string {
  const line = task.spec.trim().split('\n')[0]?.trim() ?? '';
  if (line.length > 0) {
    return line.length > 28 ? `${line.slice(0, 28)}…` : line;
  }
  return `Idea ${task.ideaId.slice(-6)}`;
}

function assigneeForTask(task: Task, agents: Agent[]): { label: string; initials: string } | null {
  if (!task.currentExecutionAgentId) return null;
  const a = agents.find((x) => x.id === task.currentExecutionAgentId);
  if (!a) {
    return { label: 'Assigned', initials: '•' };
  }
  const parts = a.name.trim().split(/\s+/).filter(Boolean);
  const initials =
    parts.length >= 2
      ? `${parts[0][0]!}${parts[parts.length - 1][0]!}`.toUpperCase()
      : (parts[0]?.slice(0, 2).toUpperCase() ?? '?');
  return { label: a.name, initials };
}

export function TaskCard({ task, agents, onOpen }: Props) {
  const assignee = assigneeForTask(task, agents);

  return (
    <article className="ft-task-card ft-animate-slide-in" style={{ display: 'flex', gap: '0.35rem', alignItems: 'flex-start' }}>
      <button
        type="button"
        className="ft-btn-icon"
        draggable
        onDragStart={(e) => {
          e.dataTransfer.setData('text/task-id', task.id);
          e.dataTransfer.effectAllowed = 'move';
          e.currentTarget.closest('.ft-task-card')?.classList.add('ft-task-card--dragging');
        }}
        onDragEnd={(e) => {
          e.currentTarget.closest('.ft-task-card')?.classList.remove('ft-task-card--dragging');
        }}
        title="Drag to another column"
        aria-label="Drag task to another column"
        style={{ flexShrink: 0, padding: '0.15rem', margin: '-0.15rem 0 0 -0.15rem' }}
      >
        <GripVertical size={14} className="ft-muted" />
      </button>
      <button
        type="button"
        onClick={() => onOpen()}
        className="ft-task-card-body-btn"
        style={{
          flex: 1,
          minWidth: 0,
          textAlign: 'left',
          background: 'none',
          border: 'none',
          padding: 0,
          color: 'inherit',
          cursor: 'pointer',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'flex-start', gap: '0.4rem' }}>
          <span className={statusDotClass(task.status)} title={task.status} aria-hidden />
          <div style={{ flex: 1, minWidth: 0 }}>
            <div className="ft-task-card-title ft-task-card-title--clamp">{task.title}</div>
            <div className="ft-task-card-tag">{taskTag(task)}</div>
            <div className="ft-task-meta ft-task-card-footer">
              {assignee ? (
                <span className="ft-task-assignee" title={assignee.label}>
                  <span className="ft-task-assignee-avatar">{assignee.initials}</span>
                  <span>{formatRelativeTime(task.updatedAt)}</span>
                </span>
              ) : (
                <span>Updated {formatRelativeTime(task.updatedAt)}</span>
              )}
            </div>
          </div>
        </div>
      </button>
    </article>
  );
}
