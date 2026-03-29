package queue

import (
	"context"

	"github.com/go-redis/redis"
)

type Storage struct {
	Queue interface {
		Add(context.Context, *PageTask) error
		PopLeft(context.Context) (*PageTask, error)
		Peek(context.Context) (*PageTask, error) // Shows us the next Crawl
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{Queue: &RedisQueue{rbd}}
}
