-- External / customer feedback linked to products (Mission Control–style).

CREATE TABLE IF NOT EXISTS product_feedback (
  id TEXT PRIMARY KEY,
  product_id TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  source TEXT NOT NULL,
  content TEXT NOT NULL,
  customer_id TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL DEFAULT '',
  sentiment TEXT NOT NULL DEFAULT 'neutral',
  processed INTEGER NOT NULL DEFAULT 0,
  idea_id TEXT REFERENCES ideas(id) ON DELETE SET NULL,
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_product_feedback_product ON product_feedback(product_id, created_at DESC);

UPDATE arms_schema_version SET version = 21 WHERE singleton = 1;
