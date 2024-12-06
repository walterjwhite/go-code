package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	google_api "github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway"
)

func (p *Provider) ReadToken(session *gateway.Session) {
	ctx, client, _ := google_api.Initialize(p.CredentialsFile, p.ProjectId, p.TokenTopicName)
	defer client.Close()

	sctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sub := client.Subscription(p.TokenSubscriptionName)

	log.Info().Msgf("subscribed to: %v", p.TokenSubscriptionName)
	err := sub.Receive(sctx, func(ctx context.Context, m *pubsub.Message) {
		m.Ack()

		log.Info().Msgf("received message: %v | %s", m.Data, m.Data)
		p.PublishStatus(fmt.Sprintf("received message: %s", m.Data), true)

		var token string
		err := json.Unmarshal(m.Data, &token)
		if err != nil {
			log.Info().Msgf("Error unmarshalling message: %v", err)
			p.PublishStatus(fmt.Sprintf("error unmarshalling message: %s", m.Data), false)
		} else {
			log.Info().Msgf("unmarshalled token: %s", token)
			p.PublishStatus(fmt.Sprintf("unmarshalled token: %s", token), true)

			if session.Run(token) {
				p.PublishStatus("session is authenticated", true)

				cancel()
				return
			}

			p.PublishStatus("failed to authenticate", false)
		}
	})

	logging.Panic(err)
}
