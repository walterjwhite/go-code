package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartDailyRequestReportWorker(t *testing.T) {
	db, err := NewSQLXDBSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL DEFAULT (datetime('now')),
			ip TEXT,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INTEGER
		)
	`)
	assert.NoError(t, err)

	_, err = db.Exec(`INSERT INTO http_requests (ts, ip, method, request_uri, user_agent, status) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		time.Now(), "127.0.0.1", "GET", "/test", "TestAgent", 200)
	assert.NoError(t, err)

	StartDailyRequestReportWorker(db)

	time.Sleep(100 * time.Millisecond)

}

func TestRunDailyReport_MissingEnvVars(t *testing.T) {
	db, err := NewSQLXDBSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL DEFAULT (datetime('now')),
			ip TEXT,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INTEGER
		)
	`)
	assert.NoError(t, err)

	_ = os.Unsetenv("REPORT_EMAIL_FROM")
	_ = os.Unsetenv("REPORT_EMAIL_TO")
	_ = os.Unsetenv("REPORT_SMTP_HOST")
	_ = os.Unsetenv("REPORT_SMTP_PORT")

	runDailyReport(db)
}

func TestRunDailyReport_EmptyTable(t *testing.T) {
	db, err := NewSQLXDBSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL DEFAULT (datetime('now')),
			ip TEXT,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INTEGER
		)
	`)
	assert.NoError(t, err)

	_ = os.Unsetenv("REPORT_EMAIL_FROM")
	_ = os.Unsetenv("REPORT_EMAIL_TO")
	_ = os.Unsetenv("REPORT_SMTP_HOST")
	_ = os.Unsetenv("REPORT_SMTP_PORT")

	runDailyReport(db)
}

func TestRequestLog_DBInsert(t *testing.T) {
	db, err := NewSQLXDBSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer func() { _ = db.Close() }()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL DEFAULT (datetime('now')),
			ip TEXT,
			method TEXT,
			request_uri TEXT,
			user_agent TEXT,
			status INTEGER
		)
	`)
	assert.NoError(t, err)

	rl := RequestLog{
		TS:         time.Now(),
		IP:         "192.168.1.1",
		Method:     "POST",
		RequestURI: "/api/users",
		UserAgent:  "TestClient/1.0",
		Status:     201,
	}

	_, err = db.Exec(`INSERT INTO http_requests (ts, ip, method, request_uri, user_agent, status) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		rl.TS, rl.IP, rl.Method, rl.RequestURI, rl.UserAgent, rl.Status)
	assert.NoError(t, err)

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM http_requests")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)
}
