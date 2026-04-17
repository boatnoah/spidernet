package main

import (
	"context"
	"log"
	"time"

	"github.com/boatnoah/spidernet/internal/crawler"
)

type worker struct {
	svc *crawler.CrawlerService
}

type config struct {
	addr     string
	db       dbConfig
	redisCfg redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (w *worker) run() error {
	log.Print("Starting worker...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer log.Print("Ending worker...")
	defer cancel()

	for {
		err := w.svc.ProcessTask(ctx)
		if err != nil {
			return err
		}
	}
}
