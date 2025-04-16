package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

func (s *Session) Subscribe(topicName string, subscriptionName string, ms MessageSubscriber) {
	_ = s.getOrCreateTopic(topicName)

	sub := s.client.Subscription(subscriptionName)
	log.Info().Msgf("subscribed to: %s", subscriptionName)

	err := sub.Receive(s.Ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Info().Msgf("received raw message: %v", m)

		decrypted := s.decrypt(m.Data)
		decompressed := s.decompress(decrypted)

		log.Info().Msgf("received message (decompressed): %s", decompressed)
		var deserialized = ms.New()
		err := json.Unmarshal(decompressed, &deserialized)
		if err != nil {
			ms.MessageParseError(err)
			logging.Panic(err)
		}

		log.Info().Msgf("received message: %s", deserialized)

		ms.MessageDeserialized()

		m.Ack() // Acknowledge that we've consumed the message.
	})

	logging.Panic(err)
}

func (s *Session) decrypt(data []byte) []byte {
	if s.AesConf == nil {
		return data
	}

	return s.AesConf.Decrypt(data)
}

func (s *Session) decompress(data []byte) []byte {
	if !s.EnableCompression {
		return data
	}

	decompressed, err := zstd.DecompressBuffer(data)
	logging.Panic(err)

	return decompressed
}

type MessageSubscriber interface {
	New() any

	MessageDeserialized()
	MessageParseError(err error)
}
