package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewSQLXDBSQLite(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, db.Ping()
}
