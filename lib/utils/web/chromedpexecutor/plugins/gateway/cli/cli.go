package cli

import (
	"flag"
	"os"
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
	if ( len(*tokenFlag) == 6) {
		return *tokenFlag
	}

	if len(os.Args) >= 2 {
		// use first argument
		return os.Args[1]
	}

	return ""
}
