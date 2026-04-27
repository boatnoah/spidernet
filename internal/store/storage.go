package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateLink     = errors.New("duplicate link")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Links interface {
		CreateBatch(ctx context.Context, jobID uuid.UUID, fromURL string, toURLs []string, depth int) error
		GetAllLinksByJobID(context.Context, uuid.UUID) ([]Links, error)
	}
	Pages interface {
		Create(context.Context, PageRequestInfo) error
	}
	CrawlJobs interface {
		CreateJob(context.Context, CrawlJobPayload) (*JobID, error)
		UpdateStatus(context.Context, uuid.UUID, string) error
		GetJobById(context.Context, uuid.UUID) (*CrawlJob, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Links:     &LinkStore{db: db},
		Pages:     &PageStore{db},
		CrawlJobs: &CrawlJobStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
