package cli

import (
	"context"

	"flag"
	"fmt"
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

func (p *Provider) OnSuccess(ctx context.Context) {
	fmt.Println("successfully authenticated")
}

func (p *Provider) OnError(ctx context.Context, err error) {
	fmt.Println("error during authentication", err)
}
