package queue

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type InMemory struct {
	Queue       []PageTask
	Outstanding map[uuid.UUID]int64
}

var ErrQueueEmpty = errors.New("Queue is empty")

func NewInMemory() *InMemory {
	return &InMemory{
		Queue:       make([]PageTask, 0, 8),
		Outstanding: make(map[uuid.UUID]int64),
	}
}

func (im *InMemory) Add(ctx context.Context, pt *PageTask) error {
	im.Outstanding[pt.CrawlJobID]++
	im.Queue = append(im.Queue, *pt)
	return nil
}
func (im *InMemory) BlockingPop(ctx context.Context) (*PageTask, error) {
	if len(im.Queue) == 0 {
		return nil, ErrQueueEmpty
	}
	first := im.Queue[0]
	im.Queue = im.Queue[1:]
	return &first, nil

}

func (im *InMemory) Peek(ctx context.Context) (*PageTask, error) {
	if len(im.Queue) == 0 {
		return nil, ErrQueueEmpty
	}
	return &im.Queue[0], nil
}

func (im *InMemory) IncrementOutstanding(ctx context.Context, jobID uuid.UUID) (int64, error) {
	im.Outstanding[jobID]++
	return im.Outstanding[jobID], nil
}

func (im *InMemory) DecrementOutstanding(ctx context.Context, jobID uuid.UUID) (int64, error) {
	im.Outstanding[jobID]--
	return im.Outstanding[jobID], nil
}
