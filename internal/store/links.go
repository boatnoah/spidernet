package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

//
// CREATE TABLE links (
// crawl_job_id UUID NOT NULL REFERENCES crawl_jobs(id) ON DELETE CASCADE,
// from_url TEXT NOT NULL,
// to_url TEXT NOT NULL,
// depth INT NOT NULL,
// PRIMARY KEY (crawl_job_id, from_url, to_url)
// );

type Links struct {
	CrawlJobID uuid.UUID `json:"crawl_job_id"`
	FromURL    string    `json:"from_url"`
	ToURL      string    `json:"to_url"`
	Depth      int       `json:"depth"`
}

type LinkStore struct {
	db *sql.DB
}

func (s *LinkStore) Create(ctx context.Context, linkPayload Links) error {
	query := `
		INSERT into links (crawl_job_id, from_url, to_url, depth) 
		VALUES ($1, $2, $3, $4)		
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		linkPayload.CrawlJobID,
		linkPayload.FromURL,
		linkPayload.ToURL,
		linkPayload.Depth,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func (s *LinkStore) GetAllLinksByJobID(ctx context.Context, jobID uuid.UUID) (*[]Links, error) {
	query := `
		SELECT * FROM links
		WHERE job_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var links []Links

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var link Links

		err = rows.Scan(&link.CrawlJobID, &link.FromURL, &link.ToURL, &link.Depth)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return &links, nil
}
