package google

import (
	"cloud.google.com/go/pubsub/v2"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

func (c *Conf) Subscribe(topicName string, subscriptionName string, ms MessageSubscriber) {
	sub := c.client.Subscriber(subscriptionName)
	log.Info().Msgf("subscribed to: %s", subscriptionName)

	err := sub.Receive(c.ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Info().Msgf("received raw message: %v", m)

		decrypted, err := c.decrypt(m.Data)
		if err != nil {
			logging.Warn(err, false, "pubsub.Subscribe.Receive")
			return
		}

		decompressed := c.decompress(decrypted)
		log.Info().Msgf("received message (decompressed): %s", decompressed)

		deserialized, err := c.deserialize(ms, decompressed)
		if err != nil {
			logging.Warn(err, false, "pubsub.Subscribe.Deserialize")
			return
		}

		ms.MessageDeserialized(deserialized)

		m.Ack() // Acknowledge that we've consumed the message.
	})
	logging.Warn(err, false, "google_pubsub.subscribe.Receive")
}

func (c *Conf) decrypt(data []byte) ([]byte, error) {
	if c.aes == nil {
		return data, nil
	}

	return c.aes.Decrypt(data)
}

func (c *Conf) decompress(data []byte) []byte {
	if !c.Compress {
		return data
	}

	decompressed, err := zstd.DecompressBuffer(data)
	logging.Warn(err, false, "google_pubsub.subscribe.decompress")

	return decompressed
}

func (c *Conf) deserialize(ms MessageSubscriber, decompressed []byte) ([]byte, error) {
	if !c.Serialize {
		log.Debug().Msg("not deserializing")
		return decompressed, nil
	}

	var deserialized []byte
	err := json.Unmarshal(decompressed, &deserialized)
	if err != nil {
		ms.MessageParseError(err)
		logging.Warn(err, false, "google_pubsub.subscribe.Unmarshal")
		return nil, err
	}

	log.Info().Msgf("received message: %s", deserialized)
	return deserialized, nil
}

type MessageSubscriber interface {
	MessageDeserialized(deserialized []byte)
	MessageParseError(err error)
}
