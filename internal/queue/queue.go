package queue

import (
	"context"

	"github.com/go-redis/redis"
)

type RedisQueue struct {
	rds *redis.Client
}

type Job struct {
	// TODO: fill out the fields
}

func (rq *RedisQueue) Add(ctx context.Context) error {
	return nil

}

func (rq *RedisQueue) PopLeft(ctx context.Context) *Job {
	return nil
}

func (rq *RedisQueue) Peek(ctx context.Context) *Job {
	return nil
}
