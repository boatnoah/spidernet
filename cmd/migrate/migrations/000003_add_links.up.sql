CREATE TABLE links (
  crawl_job_id UUID NOT NULL REFERENCES crawl_jobs (id) ON DELETE CASCADE,
  from_url TEXT NOT NULL,
  to_url TEXT NOT NULL,
  depth INT NOT NULL,
  PRIMARY KEY (crawl_job_id, from_url, to_url)
);
