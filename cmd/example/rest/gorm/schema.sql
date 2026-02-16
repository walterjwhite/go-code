-- SQL schema for storing HTTP request logs
-- Use timestamptz so times are timezone-aware
CREATE TABLE IF NOT EXISTS http_requests (
  id BIGSERIAL PRIMARY KEY,
  ts TIMESTAMPTZ NOT NULL DEFAULT now(),
  ip INET,
  method TEXT,
  request_uri TEXT,
  user_agent TEXT,
  status INT
);

-- Indexes for queries (e.g., recent queries by time, ip lookups)
CREATE INDEX IF NOT EXISTS idx_http_requests_ts ON http_requests (ts DESC);
CREATE INDEX IF NOT EXISTS idx_http_requests_ip ON http_requests (ip);
