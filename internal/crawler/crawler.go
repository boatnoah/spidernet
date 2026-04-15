package crawler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/boatnoah/spidernet/internal/queue"
	"github.com/boatnoah/spidernet/internal/store"
)

// coordinates between fetch and html go + use the storage

type CrawlerService struct {
	store  store.Storage
	queue  queue.Storage
	client *http.Client
}

func NewCrawlerService(store store.Storage, queue queue.Storage, client *http.Client) *CrawlerService {
	return &CrawlerService{
		store,
		queue,
		client,
	}
}

// rollback any of the database transacations if anything fails

func (c *CrawlerService) ProcessTask(ctx context.Context) error {
	task, err := c.queue.BlockingPop(ctx)

	if err != nil {
		return err
	}

	jobID := task.CrawlJobID
	depth := task.Depth

	job, err := c.store.CrawlJobs.GetJobById(ctx, jobID)

	if err != nil {
		return err
	}

	maxDepth := job.MaxDepth

	if depth > maxDepth {
		return nil
	}

	doc, statusCode, err := fetchURL(ctx, c.client, task.URL)

	if err != nil {
		pageInfo := store.PageRequestInfo{
			Url:        task.URL,
			Depth:      task.Depth,
			HttpStatus: statusCode,
			FetchError: fmt.Sprintf("%v", err),
		}

		c.store.Pages.Create(ctx, pageInfo)
		return err
	}

	pageInfo := store.PageRequestInfo{
		Url:        task.URL,
		Depth:      task.Depth,
		HttpStatus: statusCode,
	}

	c.store.Pages.Create(ctx, pageInfo)

	links, err := extractLinks(doc)

	if err != nil {
		return err
	}

	for _, link := range links {
		pageTask := queue.CreatePageTask(jobID, link, depth+1)

		c.queue.Add(ctx, pageTask)
	}

	return nil

}
