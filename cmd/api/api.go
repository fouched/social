package main

import (
	"github.com/fouched/social/internal/auth"
	"github.com/fouched/social/internal/mailer"
	"github.com/fouched/social/internal/repo"
	"go.uber.org/zap"
	"time"
)

type application struct {
	config        config
	repo          repo.Repository
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator //abstract authentication mechanism so that we can easily use another
}

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
