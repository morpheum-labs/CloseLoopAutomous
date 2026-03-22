-- Mission Control–style maybe pool scheduling metadata (batch re-eval / resurface).

ALTER TABLE maybe_pool ADD COLUMN last_evaluated_at TEXT NOT NULL DEFAULT '';
ALTER TABLE maybe_pool ADD COLUMN next_evaluate_at TEXT NOT NULL DEFAULT '';
ALTER TABLE maybe_pool ADD COLUMN evaluation_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE maybe_pool ADD COLUMN evaluation_notes TEXT NOT NULL DEFAULT '';

UPDATE arms_schema_version SET version = 20 WHERE singleton = 1;
