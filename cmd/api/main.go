package main

import (
	"flag"
	"github.com/fouched/social/internal/db"
	"github.com/fouched/social/internal/repo"
	"log"
)

const version = "0.0.1"

func main() {

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":9080", "Server addr to listen on")
	flag.StringVar(&cfg.env, "environment", "development", "Environment")
	flag.StringVar(&cfg.db.dsn, "dsn", "host=localhost port=5432 user=postgres password=password dbname=social sslmode=disable", "DSN (Data Source Name)")
	flag.IntVar(&cfg.db.maxOpenConn, "dbmaxconn", 10, "Max Open DB Connections")
	flag.IntVar(&cfg.db.maxIdleConn, "dbconsole", 5, "Max Idle DB Connections")

	dbPool, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConn,
		cfg.db.maxIdleConn,
	)
	if err != nil {
		log.Panic(err)
	}
	// we have database connectivity, close it after app stops
	defer dbPool.Close()
	log.Println("DB connected")

	repository := repo.NewRepository(dbPool)

	app := &application{
		config: cfg,
		repo:   repository,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
