package sqlite

import (
	"time"

	"github.com/walterjwhite/go-code/lib/net/http/request_logger"
	_ "modernc.org/sqlite"
)

const (
	DRIVER_NAME             = "sqlite"
	MAX_OPEN_CONNECTIONS    = 4
	MAX_IDLE_CONNECTIONS    = 2
	CONNECTION_MAX_LIFETIME = 30 * time.Minute
)

type SQLite struct {
}

func DefaultConf() *request_logger.Conf {
	return &request_logger.Conf{DriverName: DRIVER_NAME,
		MaxOpenConnections:    MAX_OPEN_CONNECTIONS,
		MaxIdleConnections:    MAX_IDLE_CONNECTIONS,
		ConnectionMaxLifetime: CONNECTION_MAX_LIFETIME,
	}
}

func (s *SQLite) From(c *request_logger.Conf) {
	c.DriverName = DRIVER_NAME

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = MAX_IDLE_CONNECTIONS
	}

	if c.MaxOpenConnections == 0 {
		c.MaxOpenConnections = MAX_OPEN_CONNECTIONS
	}

	if c.ConnectionMaxLifetime == 0 {
		c.ConnectionMaxLifetime = CONNECTION_MAX_LIFETIME
	}
}

func (s *SQLite) InitSQL() []string {
	return []string{
		`PRAGMA synchronous = NORMAL;`,
		`PRAGMA foreign_keys = ON;`,
		`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL DEFAULT (datetime('now')),
			ip TEXT,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_http_requests_ts ON http_requests (ts);
		CREATE INDEX IF NOT EXISTS idx_http_requests_ip ON http_requests (ip);
		`}
}
