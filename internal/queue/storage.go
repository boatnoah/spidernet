package queue

import (
	"context"

	"github.com/go-redis/redis"
)

type Storage struct {
	Queue interface {
		Add(context.Context) error
		PopLeft(context.Context) *Job
		Peek(context.Context) *Job
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{Queue: &RedisQueue{rbd}}
}
