package crawler

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (c *CrawlerService) ProcessTask(ctx context.Context) error {
	task, err := c.queue.BlockingPop(ctx)

	if err != nil {
		return err
	}

	log.Print(task)

	jobID := task.CrawlJobID
	depth := task.Depth

	var taskErr error
	defer func() {
		remaining, err := c.queue.DecrementOutstanding(ctx, jobID)
		if err != nil {
			log.Printf("unable to decrement outstanding for job %v: %v", jobID, err)
			return
		}

		if taskErr != nil {
			if err := c.store.CrawlJobs.UpdateStatus(ctx, jobID, "failed"); err != nil {
				log.Printf("unable to mark job %v failed: %v", jobID, err)
			}
			return
		}

		if remaining == 0 {
			if err := c.store.CrawlJobs.UpdateStatus(ctx, jobID, "completed"); err != nil {
				log.Printf("unable to mark job %v completed: %v", jobID, err)
			}
		}
	}()

	job, err := c.store.CrawlJobs.GetJobById(ctx, jobID)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			log.Printf("skipping job with id:%v; job not found", jobID)
			return nil
		}
		taskErr = err
		return taskErr
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
		taskErr = err
		return taskErr
	}

	pageInfo := store.PageRequestInfo{
		CrawlJobID: jobID,
		Url:        task.URL,
		Depth:      task.Depth,
		HttpStatus: statusCode,
	}

	err = c.store.Pages.Create(ctx, pageInfo)
	if err != nil {
		taskErr = err
		return taskErr
	}

	links, err := extractLinks(doc)

	if err != nil {
		taskErr = err
		return taskErr
	}

	err = c.store.Links.CreateBatch(ctx, jobID, task.URL, links, depth)
	if err != nil {
		taskErr = err
		return taskErr
	}

	for _, link := range links {
		pageTask := queue.CreatePageTask(jobID, link, depth+1)

		if err := c.queue.Add(ctx, pageTask); err != nil {
			taskErr = err
			return taskErr
		}
	}

	return nil

}
