package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boatnoah/spidernet/internal/queue"
	"github.com/boatnoah/spidernet/internal/store"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type application struct {
	config config
	store  store.Storage
	queue  queue.Storage
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

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	return nil
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}
