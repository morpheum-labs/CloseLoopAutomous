-- Dedicated preference model row per product (learning loop can evolve independently of legacy preference_model_json).

CREATE TABLE IF NOT EXISTS preference_models (
  product_id TEXT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
  model_json TEXT NOT NULL DEFAULT '{}',
  updated_at TEXT NOT NULL
);

UPDATE arms_schema_version SET version = 14 WHERE singleton = 1;
