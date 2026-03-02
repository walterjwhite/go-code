package google

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type Provider struct {
	TopicName        string
	SubscriptionName string

	Conf *google_pubsub.Conf
}

func New(ctx context.Context) (*Provider, error) {
	provider := &Provider{}
	application.Load(provider)

	if provider.Conf == nil {
		return nil, fmt.Errorf("failed to load Google Pub/Sub configuration")
	}

	if err := provider.Conf.Init(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize Google Pub/Sub: %w", err)
	}

	return provider, nil
}

func (p *Provider) String() string {
	return fmt.Sprintf("Provider: {TopicName: %s, SubscriptionName: %s, Conf: %s}", p.TopicName, p.SubscriptionName, p.Conf)
}

func (p *Provider) Publish(message []byte) error {
	return p.Conf.Publish(p.TopicName, message)
}
