package main

import (
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
}

type application struct {
	config config
	repo   repo.Repository
}

// mount initializes the router and defines it paths
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPost)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", app.getPost)
			})
		})
	})

	return r
}

// run runs the application
func (app *application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started in %s mode on port %s", app.config.env, app.config.addr)

	return srv.ListenAndServe()
}
