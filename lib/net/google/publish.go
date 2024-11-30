package google

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Publish(credentialsDirectory string, projectId string, topicName string, message interface{}) {
	ctx, client, topic := Initialize(credentialsDirectory, projectId, topicName)
	defer client.Close()

	data, err := json.Marshal(message)
	logging.Panic(err)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	id, err := result.Get(ctx)
	logging.Panic(err)

	log.Info().Msgf("Published message with ID %s, data: %v\n", id, data)
}
