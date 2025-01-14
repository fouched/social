package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/fouched/social/docs" // required for swagger docs
	"github.com/fouched/social/internal/auth"
	"github.com/fouched/social/internal/driver"
	"github.com/fouched/social/internal/mailer"
	"github.com/fouched/social/internal/repo"
	"github.com/fouched/social/internal/repo/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = "0.0.2"

//	@title			GopherSocial API
//	@description	API for GopherSocial
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":9080", "Server addr to listen on")
	flag.StringVar(&cfg.apiURL, "apiURL", "localhost:9080", "External API URL")
	flag.StringVar(&cfg.frontendURL, "frontendURL", "http://localhost:5173", "External Frontend URL")
	flag.StringVar(&cfg.env, "environment", "development", "Environment")
	flag.StringVar(&cfg.db.dsn, "dsn", "host=localhost port=5432 user=postgres password=password dbname=social sslmode=disable", "DSN (Data Source Name)")
	flag.IntVar(&cfg.db.maxOpenConn, "dbmaxconn", 10, "Max Open DB Connections")
	flag.IntVar(&cfg.db.maxIdleConn, "dbconsole", 5, "Max Idle DB Connections")
	flag.DurationVar(&cfg.mail.expiry, "expiry", time.Hour*24*3, "Registration Expiry")
	flag.StringVar(&cfg.mail.fromEmail, "fromEmail", "fouche@limehouse.co.za", "From email")
	flag.StringVar(&cfg.mail.sendgrid.apiKey, "sendgridApiKey", "", "SendGrid API Key")
	flag.StringVar(&cfg.auth.basic.user, "basicAuthUser", "admin", "Basic auth user")
	flag.StringVar(&cfg.auth.basic.pw, "basicAuthPass", "admin", "Basic auth password")
	flag.StringVar(&cfg.auth.token.secret, "tokenSecret", "example", "Token secret")
	flag.DurationVar(&cfg.auth.token.exp, "tokenExpiry", time.Hour*24*3, "Token expiry duration")
	flag.StringVar(&cfg.auth.token.issuer, "tokenIssuer", "gophersocial", "Token issuer")
	flag.StringVar(&cfg.redis.addr, "redisAddr", "localhost:6379", "Redis address")
	flag.StringVar(&cfg.redis.pw, "redisPwd", "", "Redis Password")
	flag.IntVar(&cfg.redis.db, "redisDB", 0, "Redis DB")
	flag.BoolVar(&cfg.redis.enabled, "redisEnabled", true, "Redis enabled")

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	logger.Info("Starting app...")

	// Database
	dbPool, err := driver.New(
		cfg.db.dsn,
		cfg.db.maxOpenConn,
		cfg.db.maxIdleConn,
	)
	if err != nil {
		logger.Fatal(err)
	}
	// we have database connectivity, close it after app stops
	defer dbPool.Close()
	logger.Info("DB connected")

	//cache
	var cacheDriver *redis.Client
	var cacheInstance cache.Cache
	if cfg.redis.enabled {
		cacheDriver = cache.NewRedisClient(cfg.redis.addr, cfg.redis.pw, cfg.redis.db)
		cacheInstance = cache.NewRedisCache(cacheDriver)
		logger.Info("Redis cache connected")
	}

	repository := repo.NewRepository(dbPool)

	//seed(repository)

	mailerImpl := mailer.NewSendgridClient(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.issuer,
		cfg.auth.token.issuer,
	)

	app := &application{
		config:        cfg,
		repo:          repository,
		cache:         cacheInstance,
		logger:        logger,
		mailer:        mailerImpl,
		authenticator: jwtAuthenticator,
	}

	mux := app.routes()
	logger.Fatal(app.run(mux))
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

	// do a graceful shutdown - give the server 5 secs to finish what it is doing
	// useful for docker swarm / k8s environments
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Server started", "env", app.config.env, "addr", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}

func seed(repo repo.Repository) {
	fmt.Println("Seeding database")
	Seed(repo)
}
