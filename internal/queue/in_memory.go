package queue

import (
	"context"
	"errors"
)

type InMemory struct {
	Queue []PageTask
}

var ErrQueueEmpty = errors.New("Queue is empty")

func NewInMemory() *InMemory {
	return &InMemory{Queue: make([]PageTask, 0, 8)}
}

func (im *InMemory) Add(ctx context.Context, pt *PageTask) error {
	im.Queue = append(im.Queue, *pt)
	return nil
}
func (im *InMemory) PopLeft(ctx context.Context) (*PageTask, error) {
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
