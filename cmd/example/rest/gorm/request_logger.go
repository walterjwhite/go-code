package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver usable with sqlx
	"github.com/jmoiron/sqlx"
)

type RequestLog struct {
	TS         time.Time `db:"ts"`
	IP         string    `db:"ip"`
	Method     string    `db:"method"`
	RequestURI string    `db:"request_uri"`
	UserAgent  string    `db:"user_agent"`
	Status     int       `db:"status"`
}

func NewSQLXDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, db.Ping()
}

func StartRequestLogWorker(db *sqlx.DB, bufSize int) chan<- RequestLog {
	ch := make(chan RequestLog, bufSize)

	go func() {
		isSQLite := false
		if _, err := db.Exec(`PRAGMA journal_mode = WAL;`); err == nil {
			isSQLite = true
		}

		if isSQLite {
			_, _ = db.Exec(`PRAGMA synchronous = NORMAL;`)
			_, _ = db.Exec(`PRAGMA foreign_keys = ON;`)

			_, _ = db.Exec(`
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
		`)
		} else {
			_, _ = db.Exec(`
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
		`)
		}

		insertStmt := `
	INSERT INTO http_requests (ts, ip, method, request_uri, user_agent, status)
	VALUES (:ts, :ip, :method, :request_uri, :user_agent, :status)
	`
		for rl := range ch {
			_, err := db.NamedExec(insertStmt, rl)
			if err != nil {
				log.Printf("http request log insert error: %v", err)
			}
		}
	}()

	return ch
}

func RequestLoggerMiddleware(out chan<- RequestLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		_ = start // we don't use duration here, but you could capture latency if desired

		ip := clientIP(c.Request.Header.Get("X-Forwarded-For"), c.ClientIP())
		ua := c.Request.UserAgent()
		uri := c.Request.RequestURI
		method := c.Request.Method
		rl := RequestLog{
			TS:         time.Now().UTC(),
			IP:         ip,
			Method:     method,
			RequestURI: uri,
			UserAgent:  ua,
			Status:     c.Writer.Status(),
		}

		select {
		case out <- rl:
		default:
		}
	}
}

func clientIP(forwarded string, fallback string) string {
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				if net.ParseIP(p) != nil {
					return p
				}
			}
		}
	}
	host, _, err := net.SplitHostPort(fallback)
	if err == nil {
		return host
	}
	return fallback
}
