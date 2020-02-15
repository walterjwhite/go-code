package cli

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/citrix/gateway"
)

type Provider struct {
}

var (
	tokenFlag = flag.String("Token", "", "RSA Token")
)

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Get() [] /*6*/ int {
	return gateway.Get(*tokenFlag)
}
