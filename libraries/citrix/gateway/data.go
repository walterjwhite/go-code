package gateway

import (
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/periodic"
	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Token string

	Tickle *Tickle

	chromedpsession *chromedpexecutor.ChromeDPSession
}

type Credentials struct {
	Domain   string
	Username string
	Password string

	Pin string
}

type Endpoint struct {
	Uri string

	UsernameXPath    string
	PasswordXPath    string
	TokenXPath       string
	LoginButtonXPath string
}

type Tickle struct {
	TickleInterval   time.Duration
	periodicInstance *periodic.PeriodicInstance

	lastTickledAll bool
}
