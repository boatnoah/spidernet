package main

import (
	"log"

	"github.com/boatnoah/spidernet/internal/adapter"
	"github.com/boatnoah/spidernet/internal/env"
	"github.com/boatnoah/spidernet/internal/queue"
	"github.com/boatnoah/spidernet/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),

		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/spidernet?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
	}
	db, err := adapter.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		panic(err)
	}

	store := store.NewStorage(db)

	queue := queue.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)

	app := &application{
		config: cfg,
		store:  store,
		queue:  queue,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
