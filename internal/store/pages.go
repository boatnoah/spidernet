package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Pages struct {
	CrawlJobID uuid.UUID `json:"crawl_job_id"`
	Url        string    `json:"url"`
	Depth      string    `json:"depth"`
	HttpStatus int       `json:"http_status"`
	FetchError string    `json:"fetch_error"`
}

type PageStore struct {
	db *sql.DB
}

type PageRequestInfo struct {
	CrawlJobID uuid.UUID `json:"crawl_job_id"`
	Url        string    `json:"url"`
	Depth      int       `json:"depth"`
	HttpStatus int       `json:"http_status"`
	FetchError string    `json:"fetch_error"`
}

func (s *PageStore) Create(ctx context.Context, p PageRequestInfo) error {
	query := `
		INSERT INTO pages (crawl_job_id, url, depth, http_status, fetch_error)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (crawl_job_id, url) DO NOTHING;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.CrawlJobID,
		p.Url,
		p.Depth,
		p.HttpStatus,
		p.FetchError,
	).Err()

	if err != nil {
		return err
	}
	return nil
}
