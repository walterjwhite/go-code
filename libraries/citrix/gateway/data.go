package gateway

import (
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/periodic"
	"time"
)

type TokenProvider interface {
	Get() string
}

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Token string

	Tickle *Tickle

	UseLightVersion bool

	PostAuthenticationActions []Action

	chromedpsession *chromedpexecutor.ChromeDPSession
}

type Action struct {
	Name    string
	Actions []string
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

	AuthenticationDelay *time.Duration
}

type Tickle struct {
	TickleInterval   *time.Duration
	periodicInstance *periodic.PeriodicInstance
}
