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
