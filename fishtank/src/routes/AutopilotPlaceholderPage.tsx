import { Link, useNavigate } from 'react-router-dom';
import { ChevronLeft } from 'lucide-react';

/** Reserved for cross-product autopilot; per-product flow lives under workspace Ideation (see copy below). */
export function AutopilotPlaceholderPage() {
  const navigate = useNavigate();
  return (
    <div className="ft-screen">
      <header className="ft-border-b" style={{ padding: '1rem', background: 'var(--mc-bg-secondary)' }}>
        <button type="button" className="ft-btn-ghost" onClick={() => navigate(-1)} style={{ display: 'inline-flex', alignItems: 'center', gap: '0.35rem' }}>
          <ChevronLeft size={18} />
          Back
        </button>
      </header>
      <main className="ft-container" style={{ paddingBlock: '2rem' }}>
        <h1 style={{ fontSize: '1.35rem', fontWeight: 700, marginBottom: '0.5rem' }}>Autopilot hub</h1>
        <p className="ft-muted" style={{ maxWidth: '38rem', lineHeight: 1.6, marginBottom: '1rem' }}>
          This route is a <strong style={{ color: 'inherit' }}>placeholder</strong> for a future workspace-wide autopilot view
          (multi-product summaries, global schedules, and similar). It is <strong style={{ color: 'inherit' }}>not</strong> the
          same screen as <strong style={{ color: 'inherit' }}>Ideation</strong> under a product.
        </p>
        <p style={{ maxWidth: '38rem', lineHeight: 1.6, marginBottom: '0.75rem', fontSize: '0.95rem' }}>
          <strong>Where things live today</strong>
        </p>
        <ul className="ft-muted" style={{ maxWidth: '38rem', lineHeight: 1.65, margin: '0 0 1.25rem', paddingLeft: '1.25rem' }}>
          <li style={{ marginBottom: '0.35rem' }}>
            Open a product, then <strong style={{ color: 'inherit' }}>Ideation</strong> — manual submit, SOP workshop,{' '}
            <strong style={{ color: 'inherit' }}>guided research &amp; AI drafts</strong> (stage-gated Run research / Run
            ideation).
          </li>
          <li style={{ marginBottom: '0.35rem' }}>
            <strong style={{ color: 'inherit' }}>Approvals</strong> — swipe / triage on draft ideas.
          </li>
          <li>
            <strong style={{ color: 'inherit' }}>Calendar</strong> — product autopilot cadence / schedule for that workspace.
          </li>
        </ul>
        <p className="ft-muted" style={{ margin: 0, fontSize: '0.88rem' }}>
          <Link to="/">Back to dashboard</Link> to pick a product.
        </p>
      </main>
    </div>
  );
}
