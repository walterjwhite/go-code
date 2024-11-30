package google

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

// https://medium.com/@pinkudebnath/publish-subscribe-in-google-cloud-platform-using-go-client-libraries-e8ea83b606f2
// https://pkg.go.dev/cloud.google.com/go/pubsub#hdr-Receiving
func Subscribe(credentialsDirectory string, projectId string, topicName string, subscriptionName string) {
	ctx, client, _ := Initialize(credentialsDirectory, projectId, topicName)
	defer client.Close()

	sub := client.Subscription(subscriptionName)
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Info().Msgf("received message: %v", m.Data)
		log.Info().Msgf("received message (string value): %s", m.Data)
		m.Ack() // Acknowledge that we've consumed the message.
	})
	logging.Panic(err)
}
