package request_logger

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Provider interface {
	InitSQL() []string
	From(c *Conf)
}

type Conf struct {
	DriverName string
	DSN        string

	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration

	BufferSize int
	DB         *sqlx.DB

	Provider Provider
}

func (c *Conf) PostLoad(ctx context.Context) error {
	c.Provider.From(c)

	db, err := sqlx.Open(c.DriverName, c.DSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(c.MaxOpenConnections)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.ConnectionMaxLifetime)

	c.DB = db

	c.initDB()

	return db.Ping()
}

func (c *Conf) Close() {
	if c.DB == nil {
		return
	}

	_ = c.DB.Close()
	c.DB = nil
}
