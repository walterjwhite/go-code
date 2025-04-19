package citrix

import (
	"context"
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"github.com/walterjwhite/go-code/lib/utils/worker"

	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Tickle   *Tickle
	Headless bool

	UseLightVersion bool

	Delay     *time.Duration
	Instances []Instance

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	Worker worker.Conf

	Timeout *time.Duration

	ctx    context.Context
	cancel context.CancelFunc

	keepAliveChannel <-chan time.Time
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
}

type Instance struct {
	Index  int
	Action string

	InitialActionDelay *time.Duration
	TimeBetweenActions *time.Duration

	Actions []string

	Worker CitrixWorker

	ctx    context.Context
	cancel context.CancelFunc

	session *Session

	initialized        bool
	actionsInitialized bool
}
