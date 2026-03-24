-- Last OpenClaw-class WebSocket probe outcome (pairing / policy close) for Mission Control visibility.

ALTER TABLE gateway_endpoints ADD COLUMN connection_status TEXT NOT NULL DEFAULT '';
ALTER TABLE gateway_endpoints ADD COLUMN pairing_request_id TEXT NOT NULL DEFAULT '';
ALTER TABLE gateway_endpoints ADD COLUMN pairing_message TEXT NOT NULL DEFAULT '';
ALTER TABLE gateway_endpoints ADD COLUMN last_close_code INTEGER NOT NULL DEFAULT 0;

UPDATE arms_schema_version SET version = 32 WHERE singleton = 1;
