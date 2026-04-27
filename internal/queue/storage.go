package queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Add(context.Context, *PageTask) error
	BlockingPop(context.Context) (*PageTask, error)
	Peek(context.Context) (*PageTask, error) // Shows us the next Crawl
	IncrementOutstanding(context.Context, uuid.UUID) (int64, error)
	DecrementOutstanding(context.Context, uuid.UUID) (int64, error)
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return &RedisQueue{rbd}
}
