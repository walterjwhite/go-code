package pf

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	// "github.com/coredns/coredns/request"
	// "github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/plugin/pkg/dnstest"

	"github.com/miekg/dns"
	// "net"
	// "strconv"
	"fmt"
	// "time"
)

const pluginName = "pf"
type Pf struct {
	Next plugin.Handler
}

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	// no args after pf
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

	// tpe, _ := response.Typify(rrw.Msg, time.Now().UTC())
	// class := response.Classify(tpe)

	fmt.Printf("pf - ServeDNS - %v\n", rrw.Msg)
	return rc, err
}

func (p Pf) Name() string { return pluginName }
