package pf

import (
	"context"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"

	"github.com/miekg/dns"
	"fmt"
)

const pluginName = "pf"

type Pf struct {
	Next plugin.Handler
}

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error(pluginName, c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Pf{Next: next}
	})

	return nil
}

func (p Pf) Ready() bool {
	return true
}

func (p Pf) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	rrw := dnstest.NewRecorder(w)
	rc, err := plugin.NextOrFailure(p.Name(), p.Next, ctx, rrw, r)


	fmt.Printf("pf - ServeDNS - %v\n", rrw.Msg)
	ip := "8.8.8.8"
	add(ip)

	return rc, err
}

func (p Pf) Name() string { return pluginName }
