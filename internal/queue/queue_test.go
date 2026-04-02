package queue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryQueue(t *testing.T) {
	assert := assert.New(t)

	inmemory := NewInMemory()

	store := Storage{inmemory}

	pageTasks := [...]PageTask{
		{
			CrawlJobID: "1",
			URL:        "boatnoah.com",
			Depth:      1,
		},
		{
			CrawlJobID: "1",
			URL:        "boatnoah.com/home",
			Depth:      1,
		},
		{
			CrawlJobID: "1",
			URL:        "boatnoah.com/about",
			Depth:      1,
		},
	}

	for _, pageTask := range pageTasks {
		store.Queue.Add(context.Background(), &pageTask)
	}

	assert.Len(inmemory.Queue, 3)

	firstTask, err := store.Queue.PopLeft(context.Background())

	if err != nil {
		assert.Fail("%v", err)
	}

	assert.Equal("1", firstTask.CrawlJobID)
	assert.Equal("boatnoah.com", firstTask.URL)
	assert.Equal(1, firstTask.Depth)

	nextTask, err := store.Queue.Peek(context.Background())
	assert.Equal("1", nextTask.CrawlJobID)
	assert.Equal("boatnoah.com/home", nextTask.URL)
	assert.Equal(1, nextTask.Depth)

	assert.Len(inmemory.Queue, 2)
}
