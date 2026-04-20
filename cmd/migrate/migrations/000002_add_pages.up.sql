CREATE TABLE pages (
  crawl_job_id UUID NOT NULL REFERENCES crawl_jobs (id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  depth INT NOT NULL,
  http_status INT,
  fetch_error TEXT,
  PRIMARY KEY (crawl_job_id, url)
);
