package cli

import (
	"flag"
)

type Provider struct {
}

var (
	tokenFlag = flag.String("t", "", "RSA Token")
)

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Get() string {
	return *tokenFlag
}
