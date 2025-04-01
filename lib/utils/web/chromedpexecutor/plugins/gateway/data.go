package gateway

import (
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Tickle        *Tickle
	Headless      bool

	UseLightVersion bool

	Delay   *time.Duration
	Instances []Instance

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	session session.ChromeDPSession
}

func (s *Session) Session() session.ChromeDPSession {
	return s.session
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

type Instance struct {
	Index int
	WiggleMouse bool
}
