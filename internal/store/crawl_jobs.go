package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// id UUID PRIMARY KEY,
// start_url TEXT NOT NULL,
// status TEXT NOT NULL, -- 'running' | 'completed' | 'failed'
// max_depth INT NOT NULL,
// created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
// completed_at TIMESTAMPTZ

type CrawlJob struct {
	ID          uuid.UUID `json:"id"`
	StartUrl    string    `json:"start_url"`
	Status      string    `json:"status"`
	MaxDepth    int       `json:"max_depth"`
	CreatedAt   string    `json:"created_at"`
	CompletedAt *string   `json:"completed_at"`
}

type CrawlJobStore struct {
	db *sql.DB
}

type CrawlJobPayload struct {
	StartUrl string
	Status   string
	MaxDepth int
}

type JobID struct {
	ID uuid.UUID
}

func (s *CrawlJobStore) CreateJob(ctx context.Context, cj CrawlJobPayload) (*JobID, error) {

	query := `INSERT INTO crawl_jobs (start_url, status, max_depth) VALUES ($1, $2, $3) RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var jobID JobID

	err := s.db.QueryRowContext(
		ctx,
		query,
		cj.StartUrl,
		cj.Status,
		cj.MaxDepth,
	).Scan(
		&jobID.ID,
	)
	if err != nil {
		return nil, err

	}

	return &jobID, nil

}

func (s *CrawlJobStore) UpdateStatus(ctx context.Context, jobID uuid.UUID, status string) error {
	query := `
		UPDATE crawl_jobs
		SET status = $1, completed_at = NOW()
		WHERE id = $2
	`

	err := s.db.QueryRowContext(ctx, query, status, jobID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (s *CrawlJobStore) GetJobById(ctx context.Context, jobID uuid.UUID) (*CrawlJob, error) {
	query := `
		SELECT * from crawl_jobs
		WHERE id = $1
	`

	var job CrawlJob

	err := s.db.QueryRowContext(ctx, query, jobID).Scan(
		&job.ID,
		&job.StartUrl,
		&job.Status,
		&job.MaxDepth,
		&job.CreatedAt,
		&job.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}
