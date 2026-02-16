package google

import (
	"context"
	"fmt"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type Provider struct {
	Conf *google_pubsub.Conf

	TokenTopicName        string
	TokenSubscriptionName string

	StatusTopicName        string
	StatusSubscriptionName string

	token string
}

func (p *Provider) String() string {
	return fmt.Sprintf("Provider: {TokenTopicName: %s, TokenSubscriptionName: %s, StatusTopicName: %s, StatusSubscriptionName: %s, Conf: %s}", p.TokenTopicName, p.TokenSubscriptionName, p.StatusTopicName,
		p.StatusSubscriptionName, p.Conf)
}

func (p *Provider) Init(ctx context.Context) {
	p.Conf.Init(ctx)
}

func (p *Provider) Cleanup() {
	p.Conf.Cancel()
}
