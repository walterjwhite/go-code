package request_logger

import (
	"time"
)

type RequestLog struct {
	TS         time.Time `db:"ts"`
	IP         string    `db:"ip"`
	Method     string    `db:"method"`
	RequestURI string    `db:"request_uri"`
	UserAgent  string    `db:"user_agent"`
	Status     int       `db:"status"`
}
