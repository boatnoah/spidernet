package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Add(context.Context, *PageTask) error
	BlockingPop(context.Context) (*PageTask, error)
	Peek(context.Context) (*PageTask, error) // Shows us the next Crawl
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return &RedisQueue{rbd}
}
