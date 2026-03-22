-- Soft-delete for products (Mission Control parity); NULL = active.
ALTER TABLE products ADD COLUMN deleted_at TEXT;

UPDATE arms_schema_version SET version = 25 WHERE singleton = 1;
