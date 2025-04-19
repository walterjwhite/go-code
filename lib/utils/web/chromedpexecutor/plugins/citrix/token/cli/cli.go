package cli

import (
	"context"

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

func (p *Provider) ReadToken(ctx context.Context) *string {
	if len(*tokenFlag) == 6 {
		return tokenFlag
	}

	return nil
}


