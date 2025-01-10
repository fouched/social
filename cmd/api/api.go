package main

import (
	"fmt"
	"github.com/fouched/social/docs" // required for swagger docs
	"github.com/fouched/social/internal/auth"
	"github.com/fouched/social/internal/mailer"
	"github.com/fouched/social/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type config struct {
	addr        string
	db          dbConfig
	mail        mailConfig
	auth        authConfig
	env         string
	apiURL      string
	mailer      mailer.Client
	frontendURL string
}

type mailConfig struct {
	sendgrid  sendgridConfig
	mailtrap  mailtrapConfig
	fromEmail string
	expiry    time.Duration
}

type sendgridConfig struct {
	apiKey string
}

type mailtrapConfig struct {
	apiKey string
}

type dbConfig struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	issuer string
}

type application struct {
	config        config
	repo          repo.Repository
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator //abstract authentication mechanism so that we can easily use another
}

// mount initializes the router and defines it paths
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUser)
			r.Post("/token", app.createToken)
		})

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPost)

			r.Route("/{id}", func(r chi.Router) {
				// use middleware to retrieve post for all post handlers
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPost)
				r.Patch("/", app.updatePost)
				r.Delete("/", app.deletePost)

				r.Post("/comment", app.createPostComment)
			})
		})

		r.Route("/users", func(r chi.Router) {

			r.Put("/activate/{token}", app.activateUser)

			r.Route("/{id}", func(r chi.Router) {
				// use middleware to retrieve user for all user handlers
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUser)

				r.Put("/follow", app.followUser)
				r.Put("/unfollow", app.unfollowUser)
			})

			//Group adds a new inline-Router along the current routing path,
			//with a fresh middleware stack for the inline-Router
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeed)
			})
		})

	})

	return r
}

// run runs the application
func (app *application) run(mux http.Handler) error {
	//Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server started", "env", app.config.env, "addr", app.config.addr)

	return srv.ListenAndServe()
}
