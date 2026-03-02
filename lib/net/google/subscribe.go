package google

import (
	"cloud.google.com/go/pubsub/v2"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

func (c *Conf) processMessage(ms MessageSubscriber, data []byte) error {
	decrypted, err := c.decrypt(data)
	if err != nil {
		return err
	}

	decompressed, err := c.decompress(decrypted)
	if err != nil {
		return err
	}

	deserialized, err := c.deserialize(ms, decompressed)
	if err != nil {
		return err
	}

	ms.MessageDeserialized(deserialized)
	return nil
}

func (c *Conf) Subscribe(topicName string, subscriptionName string, ms MessageSubscriber) {
	sub := c.client.Subscriber(subscriptionName)
	log.Info().Msgf("subscribed to: %s", subscriptionName)

	err := sub.Receive(c.ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Debug().Msg("received message from subscription")

		err := c.processMessage(ms, m.Data)
		if err != nil {
			logging.Warn(err, "pubsub.Subscribe.ProcessMessage")
		}

		m.Ack() // Acknowledge that we've consumed the message.
	})
	logging.Warn(err, "google_pubsub.subscribe.Receive")
}

func (c *Conf) decrypt(data []byte) ([]byte, error) {
	if c.encryptor == nil {
		return data, nil
	}

	return c.encryptor.Decrypt(data)
}

func (c *Conf) decompress(data []byte) ([]byte, error) {
	if !c.Compress {
		return data, nil
	}

	decompressed, err := zstd.DecompressBuffer(data)
	if err != nil {
		logging.Warn(err, "google_pubsub.subscribe.decompress")
		return nil, err
	}

	return decompressed, nil
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
		logging.Warn(err, "google_pubsub.subscribe.Unmarshal")
		return nil, err
	}

	log.Debug().Msg("message deserialized successfully")
	return deserialized, nil
}

type MessageSubscriber interface {
	MessageDeserialized(deserialized []byte)
	MessageParseError(err error)
}
