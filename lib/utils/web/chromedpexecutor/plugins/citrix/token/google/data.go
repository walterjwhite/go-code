package google

import (
	"context"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type Provider struct {
	Conf *google_pubsub.Conf

	TokenTopicName        string
	TokenSubscriptionName string

	StatusTopicName        string
	StatusSubscriptionName string

	session *google_pubsub.Session
	token   string
}

func (p *Provider) Init(ctx context.Context) {
	p.session = google_pubsub.New(p.Conf.CredentialsFile, p.Conf.ProjectId, ctx)
}
