package citrix

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/token/providers/google"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
	"github.com/walterjwhite/go-code/lib/utils/worker"

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

	Worker worker.Conf

	Timeout *time.Duration

	ProxyServerAddress string

	GoogleProvider *google.Provider

	ctx    context.Context
	cancel context.CancelFunc

	controller *chromedpexecutor.ChromeDPController
}

func (s *Session) String() string {
	return fmt.Sprintf("Endpoint: %s", s.Endpoint)
}

func (s *Session) Close() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.Credentials != nil {
		s.Credentials.Clear()
	}
}





type Credentials struct {
	Domain   string
	Username string
	Password string

	Pin string
}

func (c *Credentials) Clear() {
	for i := range c.Password {
		c.Password = c.Password[:i] + "\x00" + c.Password[i+1:]
	}
	c.Password = ""

	for i := range c.Pin {
		c.Pin = c.Pin[:i] + "\x00" + c.Pin[i+1:]
	}
	c.Pin = ""

	c.Username = ""
	c.Domain = ""
}

type Endpoint struct {
	Uri string

	AuthenticationDelay *time.Duration
}

type Tickle struct {
	TickleInterval *time.Duration
}
