package main

import (
	"github.com/fouched/social/docs" // required for swagger docs
	"github.com/fouched/social/internal/auth"
	"github.com/fouched/social/internal/mailer"
	"github.com/fouched/social/internal/repo"
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
