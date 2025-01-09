package main

import (
	"flag"
	"fmt"
	"github.com/fouched/social/internal/db"
	"github.com/fouched/social/internal/mailer"
	"github.com/fouched/social/internal/repo"
	"go.uber.org/zap"
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
	flag.StringVar(&cfg.frontendURL, "frontendURL", "http://localhost:4000", "External Frontend URL")
	flag.StringVar(&cfg.env, "environment", "development", "Environment")
	flag.StringVar(&cfg.db.dsn, "dsn", "host=localhost port=5432 user=postgres password=password dbname=social sslmode=disable", "DSN (Data Source Name)")
	flag.IntVar(&cfg.db.maxOpenConn, "dbmaxconn", 10, "Max Open DB Connections")
	flag.IntVar(&cfg.db.maxIdleConn, "dbconsole", 5, "Max Idle DB Connections")
	flag.DurationVar(&cfg.mail.expiry, "expiry", time.Hour*24*3, "Registration Expiry")
	flag.StringVar(&cfg.mail.fromEmail, "fromEmail", "fouche@limehouse.co.za", "From email")
	flag.StringVar(&cfg.mail.sendgrid.apiKey, "sendgridApiKey", "", "SendGrid API Key")

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	logger.Info("Starting app...")

	// Database
	dbPool, err := db.New(
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

	repository := repo.NewRepository(dbPool)

	//seed(repository)

	mailerImpl := mailer.NewSendgridClient(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		repo:   repository,
		logger: logger,
		mailer: mailerImpl,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}

func seed(repo repo.Repository) {
	fmt.Println("Seeding database")
	Seed(repo)
}
