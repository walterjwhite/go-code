package google

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"google.golang.org/api/option"
)

func Initialize(credentialsFile string, projectId string, topicName string) (context.Context, *pubsub.Client, *pubsub.Topic) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsFile(credentialsFile))

	logging.Panic(err)

	topic := client.Topic(topicName)
	ok, err := topic.Exists(ctx)
	logging.Panic(err)
	if !ok {
		_, err := client.CreateTopic(ctx, topicName)
		logging.Panic(err)
	}

	return ctx, client, topic
}
