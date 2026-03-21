-- Inter-subtask mail for convoys (baseline; not agent execution mailbox).

CREATE TABLE IF NOT EXISTS convoy_mail (
  id TEXT PRIMARY KEY,
  convoy_id TEXT NOT NULL REFERENCES convoys(id) ON DELETE CASCADE,
  subtask_id TEXT NOT NULL,
  body TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_convoy_mail_convoy ON convoy_mail(convoy_id, created_at DESC);

UPDATE arms_schema_version SET version = 16 WHERE singleton = 1;
