package gateway

import (
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
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

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	ChromeDPSession *chromedpexecutor.ChromeDPSession
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
