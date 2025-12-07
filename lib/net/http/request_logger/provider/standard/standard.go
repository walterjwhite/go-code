package standard

import (
	"github.com/walterjwhite/go-code/lib/net/http/request_logger"
	"time"
)

const (
	MAX_OPEN_CONNECTIONS    = 50
	MAX_IDLE_CONNECTIONS    = 10
	CONNECTION_MAX_LIFETIME = 30 * time.Minute
)

type SQLx struct {
}

func DefaultConf(driverName string) *request_logger.Conf {
	return &request_logger.Conf{DriverName: driverName,
		MaxOpenConnections:    MAX_OPEN_CONNECTIONS,
		MaxIdleConnections:    MAX_IDLE_CONNECTIONS,
		ConnectionMaxLifetime: CONNECTION_MAX_LIFETIME,
	}
}

func From(c *request_logger.Conf) {
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

func (s *SQLx) InitSQL() []string {
	return []string{
		`
	CREATE TABLE IF NOT EXISTS http_requests (
			id BIGSERIAL PRIMARY KEY,
			ts TIMESTAMPTZ NOT NULL DEFAULT now(),
			ip INET,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INT
		);
		CREATE INDEX IF NOT EXISTS idx_http_requests_ts ON http_requests (ts DESC);
		CREATE INDEX IF NOT EXISTS idx_http_requests_ip ON http_requests (ip);
	`}
}
