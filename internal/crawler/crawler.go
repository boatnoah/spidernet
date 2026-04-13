package crawler

import (
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
