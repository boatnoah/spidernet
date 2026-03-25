package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Everything is working"))
		})
		r.Post("/register", app.registerUserHandler)
		r.Post("/login", app.loginUserHandler)
		r.Route("/ranked", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/top", app.topPlayersHandler)
			r.Get("/leaderboard", app.leaderboardHandler)
			r.Post("/score", app.matchSubmissionHandler)
		})
	})

	return r
}
