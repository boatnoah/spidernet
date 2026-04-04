# Spidernet

A distributed webcrawler that recursively finds links 🕸️

# Architecture

This project uses the job-worker paradigm where each unit of work is to find links on given URL.

# Models

```Go
type CrawlJob struct {
	ID        int64
	StartURL  string
	Status    string
	MaxDepth  int
	CreatedAt time.Time
}

type PageTask struct {
	CrawlJobID string `json:"crawl_job_id"`
	URL        string `json:"url"`
	Depth      int    `json:"depth"`
}

type Page struct {
	CrawlJobID string
	URL        string
	Depth      int
	HTTPStatus *int
	FetchError *string
}

type Link struct {
	CrawlJobID string
	FromURL    string
	ToURL      string
	Depth      int
}
```

```sql
CREATE TABLE crawl_jobs (
id UUID PRIMARY KEY,
start_url TEXT NOT NULL,
status TEXT NOT NULL, -- 'running' | 'completed' | 'failed'
max_depth INT NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
completed_at TIMESTAMPTZ
);

CREATE TABLE pages (
crawl_job_id UUID NOT NULL REFERENCES crawl_jobs(id) ON DELETE CASCADE,
url TEXT NOT NULL,
depth INT NOT NULL,
http_status INT,
fetch_error TEXT,
PRIMARY KEY (crawl_job_id, url)
);

CREATE TABLE links (
crawl_job_id UUID NOT NULL REFERENCES crawl_jobs(id) ON DELETE CASCADE,
from_url TEXT NOT NULL,
to_url TEXT NOT NULL,
depth INT NOT NULL,
PRIMARY KEY (crawl_job_id, from_url, to_url)
);

```
