package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (s *LinkStore) CreateBatch(ctx context.Context, jobID uuid.UUID, fromURL string, toURLs []string, depth int) error {
	if len(toURLs) == 0 {
		return nil
	}

	args := make([]any, 0, len(toURLs)*4)
	placeholders := make([]string, 0, len(toURLs))
	for i, toURL := range toURLs {
		base := i*4 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", base, base+1, base+2, base+3))
		args = append(args, jobID, fromURL, toURL, depth)
	}

	query := fmt.Sprintf(`
		INSERT INTO links (crawl_job_id, from_url, to_url, depth)
		VALUES %s
		ON CONFLICT (crawl_job_id, from_url, to_url) DO NOTHING;
	`, strings.Join(placeholders, ","))

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *LinkStore) GetAllLinksByJobID(ctx context.Context, jobID uuid.UUID) (*[]Links, error) {
	query := `
		SELECT crawl_job_id, from_url, to_url, depth
		FROM links
		WHERE crawl_job_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var links []Links

	rows, err := s.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link Links

		err = rows.Scan(&link.CrawlJobID, &link.FromURL, &link.ToURL, &link.Depth)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &links, nil
}
