package google

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

func (s *Session) Publish(topicName string, message interface{}) {
	topic := s.getOrCreateTopic(topicName)

	data, err := json.Marshal(message)
	logging.Panic(err)

	compressed := s.compress(data)
	encrypted := s.encrypt(compressed)

	result := topic.Publish(s.Ctx, &pubsub.Message{
		Data: encrypted,
	})
	id, err := result.Get(s.Ctx)
	logging.Panic(err)

	log.Info().Msgf("published message with ID %s, message: %s, data: %s", id, message, data)
}

func (s *Session) encrypt(data []byte) []byte {
	if s.AesConf == nil {
		return data
	}

	return s.AesConf.Encrypt(data)
}

func (s *Session) compress(data []byte) []byte {
	if !s.EnableCompression {
		return data
	}

	return zstd.CompressBuffer(data)
}
