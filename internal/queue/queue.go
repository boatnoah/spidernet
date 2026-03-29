package queue

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
)

type RedisQueue struct {
	rds *redis.Client
}

type PageTask struct {
	CrawlJobID string `json:"crawl_job_id"`
	URL        string `json:"url"`
	Depth      int    `json:"depth"`
}

var queueKey = "crawlqueue"

func (rq *RedisQueue) Add(ctx context.Context, pt *PageTask) error {

	ptString, err := json.Marshal(pt)

	if err != nil {
		return err
	}

	err = rq.rds.RPush(queueKey, ptString).Err()

	if err != nil {
		return err
	}
	return nil
}

func (rq *RedisQueue) PopLeft(ctx context.Context) (*PageTask, error) {
	result, err := rq.rds.LPop(queueKey).Result()
	if err != nil {
		return nil, err
	}

	var pt PageTask

	err = json.Unmarshal([]byte(result), &pt)

	if err != nil {
		return nil, err
	}

	return &pt, nil
}

// Im on vacation lol so lets make a quick video today
// i got something wrong

func (rq *RedisQueue) Peek(ctx context.Context) (*PageTask, error) {
	result, err := rq.rds.LIndex(queueKey, 0).Result()
	if err != nil {
		return nil, err
	}

	var pt PageTask

	err = json.Unmarshal([]byte(result), &pt)

	if err != nil {
		return nil, err
	}
	return &pt, nil
}
