package queue

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryQueue(t *testing.T) {
	assert := assert.New(t)

	inmemory := NewInMemory()

	jobID := uuid.New()

	pageTasks := [...]PageTask{
		{CrawlJobID: jobID, URL: "boatnoah.com", Depth: 1},
		{CrawlJobID: jobID, URL: "boatnoah.com/home", Depth: 1},
		{CrawlJobID: jobID, URL: "boatnoah.com/about", Depth: 1},
	}

	for _, pageTask := range pageTasks {
		err := inmemory.Add(context.Background(), &pageTask)
		assert.NoError(err)
	}

	assert.Len(inmemory.Queue, 3)
	assert.Equal(int64(3), inmemory.Outstanding[jobID])

	firstTask, err := inmemory.BlockingPop(context.Background())

	if err != nil {
		assert.Fail("%v", err)
	}

	assert.Equal(jobID, firstTask.CrawlJobID)
	assert.Equal("boatnoah.com", firstTask.URL)
	assert.Equal(1, firstTask.Depth)

	nextTask, err := inmemory.Peek(context.Background())
	assert.NoError(err)
	assert.Equal(jobID, nextTask.CrawlJobID)
	assert.Equal("boatnoah.com/home", nextTask.URL)
	assert.Equal(1, nextTask.Depth)

	assert.Len(inmemory.Queue, 2)

	remaining, err := inmemory.DecrementOutstanding(context.Background(), jobID)
	assert.NoError(err)
	assert.Equal(int64(2), remaining)
}
