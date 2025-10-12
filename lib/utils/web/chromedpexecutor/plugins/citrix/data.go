package citrix

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/google"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
	"github.com/walterjwhite/go-code/lib/utils/worker"
	"sync"
	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Tickle           *Tickle
	KeepAliveTimeout *time.Duration
	KeepAliveTries   uint
	KeepAliveDelay   *time.Duration

	Conf *provider.Conf

	UseLightVersion bool

	Delay     *time.Duration
	Instances []*Instance

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	Worker worker.Conf

	Timeout *time.Duration

	ProxyServerAddress string

	GoogleProvider *google.Provider

	ctx    context.Context
	cancel context.CancelFunc

	keepAliveTicker *time.Ticker

	waitGroup *sync.WaitGroup
}

func (s *Session) String() string {
	return fmt.Sprintf("Endpoint: %s", s.Endpoint)
}

type Credentials struct {
	Domain   string
	Username string
	Password string

	Pin string
}

type Endpoint struct {
	Uri string

	AuthenticationDelay *time.Duration
}

type Tickle struct {
	TickleInterval *time.Duration
}
