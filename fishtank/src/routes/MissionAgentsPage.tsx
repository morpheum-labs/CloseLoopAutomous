import { AgentsPanel } from '../components/workspace/AgentsPanel';

/** Full-width agents view — `AgentsPanel` syncs `GET /api/agents` + `GET /api/fleet/identities` on mount; registry + task heartbeats + gateway identities. */
export function MissionAgentsPage() {
  return (
    <div className="ft-queue-flex" style={{ flex: 1, minWidth: 0, minHeight: 0, padding: '0.75rem', overflow: 'auto' }}>
      <AgentsPanel />
    </div>
  );
}
