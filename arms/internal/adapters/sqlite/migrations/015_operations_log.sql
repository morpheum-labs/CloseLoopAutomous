-- Append-only operator audit trail.

CREATE TABLE IF NOT EXISTS operations_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL,
  actor TEXT NOT NULL DEFAULT '',
  action TEXT NOT NULL DEFAULT '',
  resource_type TEXT NOT NULL DEFAULT '',
  resource_id TEXT NOT NULL DEFAULT '',
  detail_json TEXT NOT NULL DEFAULT '{}',
  product_id TEXT REFERENCES products(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_operations_log_created ON operations_log(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operations_log_product ON operations_log(product_id, created_at DESC);

UPDATE arms_schema_version SET version = 15 WHERE singleton = 1;
