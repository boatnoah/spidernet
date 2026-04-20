CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE crawl_jobs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  start_url TEXT NOT NULL,
  status TEXT NOT NULL, -- 'running' | 'completed' | 'failed'
  max_depth INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
  completed_at TIMESTAMPTZ
);
