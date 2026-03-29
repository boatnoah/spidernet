package queue

import (
	"context"
	"errors"
)

type InMemory struct {
	queue []PageTask
}

var ErrQueueEmpty = errors.New("Queue is empty")

func NewInMemory() *InMemory {
	return &InMemory{queue: make([]PageTask, 8)}
}

func (im *InMemory) Add(ctx context.Context, pt *PageTask) error {
	im.queue = append(im.queue, *pt)
	return nil
}
func (im *InMemory) PopLeft(ctx context.Context) (*PageTask, error) {
	if len(im.queue) == 0 {
		return nil, ErrQueueEmpty
	}
	first := im.queue[0]
	im.queue = im.queue[:1]
	return &first, nil

}

func (im *InMemory) Peek(ctx context.Context) (*PageTask, error) {
	if len(im.queue) == 0 {
		return nil, ErrQueueEmpty
	}
	return &im.queue[0], nil
}
