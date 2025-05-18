package citrix

import (
	"context"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
	"github.com/walterjwhite/go-code/lib/utils/worker"

	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Tickle           *Tickle
	KeepAliveTimeout *time.Duration

	Conf *provider.Conf

	UseLightVersion bool

	Delay     *time.Duration
	Instances []*Instance

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	Worker worker.Conf

	Timeout *time.Duration

	ProxyServerAddress string

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
	TickleInterval *time.Duration
}

type Instance struct {
	Index      int
	WorkerType WorkerType

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

type WorkerType int

const (
	MouseWiggler WorkerType = iota
	NOOP
)

func (w WorkerType) String() string {
	return [...]string{"MouseWiggler", "NOOP"}[w]
}
