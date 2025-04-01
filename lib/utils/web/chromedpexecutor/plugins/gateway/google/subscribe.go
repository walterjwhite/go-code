package google

import (
	"context"
	"fmt"
)

func (p *Provider) ReadToken(ctx context.Context) *string {
	p.session.Subscribe(p.TokenTopicName, p.TokenSubscriptionName, p)
	return &p.token
}

func (p *Provider) New() any {
	return &p.token
}

func (p *Provider) MessageDeserialized() {
	p.publishStatus(fmt.Sprintf("unmarshalled token: %s", p.token), true)

	p.session.Cancel()
}

func (p *Provider) MessageParseError(err error) {
	p.publishStatus(fmt.Sprintf("error unmarshalling message: %s", err), false)
}

func (p *Provider) OnSuccess(ctx context.Context) {
	p.publishStatus("session is authenticated", true)
}

func (p *Provider) OnError(ctx context.Context, err error) {
	p.publishStatus("failed to authenticated", false)
}
