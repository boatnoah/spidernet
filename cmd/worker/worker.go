package main

import (
	"context"
	"log"

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
	defer log.Print("Ending worker...")

	for {
		err := w.svc.ProcessTask(context.Background())
		if err != nil {
			return err
		}
	}
}
