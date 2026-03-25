package queue

import (
	"github.com/go-redis/redis"
)

type Storage struct {
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{}
}
