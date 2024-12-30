package main

import (
	"flag"
	"github.com/fouched/social/internal/db"
	"github.com/fouched/social/internal/store"
	"log"
)

func main() {

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":9080", "Server addr to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", "host=localhost port=5432 user=postgres password=password dbname=social sslmode=disable", "DSN (Data Source Name)")
	flag.IntVar(&cfg.db.maxOpenConn, "dbmaxconn", 10, "Max Open Connections")
	flag.IntVar(&cfg.db.maxIdleConn, "dbidleconn", 5, "Max Idle Connections")

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
	log.Print("DB connected")

	storage := store.NewStorage(dbPool)

	app := &application{
		config:  cfg,
		storage: storage,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
