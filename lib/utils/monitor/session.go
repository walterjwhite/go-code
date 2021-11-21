package monitor

import (
	"context"
	"github.com/walterjwhite/go-code/lib/time/after"
	"time"
)

type Session struct {
	Directory   string
	Name        string
	Description string

	Icon string

	Actions []Action
	Context context.Context
	Channel chan *NotificationEvent

	Alerts            []Alert
	NoActivity        NoActivity
	LastAlertDateTime time.Time
}

type Action struct {
	Interval    string
	Description string
	Type        string
	Reference   string
	Monitor     Monitor
	Session     *Session
}

type Alert struct {
	Name      string
	Type      string
	Reference string
}

type NoActivity struct {
	Interval    string
	Description string
	After       *after.AfterDelay
}
