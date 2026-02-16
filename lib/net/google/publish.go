package google

import (
	"cloud.google.com/go/pubsub/v2"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

func (c *Conf) prepareMessage(message []byte) ([]byte, error) {
	data, err := c.serialize(message)
	if err != nil {
		return nil, err
	}

	compressed := c.compress(data)
	encrypted := c.encrypt(compressed)

	return encrypted, nil
}

func (c *Conf) Publish(topicName string, message []byte) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
	}

	log.Info().Msgf("publishing to: %v", topicName)

	encrypted, err := c.prepareMessage(message)
	if err != nil {
		return err
	}

	publisher := c.client.Publisher(topicName)
	defer publisher.Stop()

	publishResult := publisher.Publish(c.ctx, &pubsub.Message{Data: encrypted})
	id, err := publishResult.Get(c.ctx)
	if err != nil {
		return err
	}

	log.Info().Msgf("published message with ID %s, message: %s", id, message)
	return nil
}

func (c *Conf) serialize(message []byte) ([]byte, error) {
	if !c.Serialize {
		return message, nil
	}

	return json.Marshal(message)
}

func (c *Conf) encrypt(data []byte) []byte {
	if c.aes == nil {
		return data
	}

	return c.aes.Encrypt(data)
}

func (c *Conf) compress(data []byte) []byte {
	if !c.Compress {
		return data
	}

	return zstd.CompressBuffer(data)
}
