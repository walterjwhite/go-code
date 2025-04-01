package google

import (
	"cloud.google.com/go/pubsub"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *Session) getOrCreateTopic(topicName string) *pubsub.Topic {
	topic := s.client.Topic(topicName)
	ok, err := topic.Exists(s.Ctx)
	logging.Panic(err)
	if !ok {
		_, err := s.client.CreateTopic(s.Ctx, topicName)
		logging.Panic(err)
	}

	return topic
}
