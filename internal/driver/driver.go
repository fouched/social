package driver

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(dsn string, maxOpenConn, maxIdleConn int) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	// test database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
