package google

import (
	"cloud.google.com/go/pubsub/v2"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
)

const MaxMessageSize = 10 * 1024 * 1024

func (c *Conf) prepareMessage(message []byte) ([]byte, error) {
	if len(message) > MaxMessageSize {
		return nil, fmt.Errorf("message size %d exceeds maximum allowed size %d", len(message), MaxMessageSize)
	}

	data, err := c.serialize(message)
	if err != nil {
		return nil, err
	}

	if len(data) > MaxMessageSize {
		return nil, errors.New("serialized message size exceeds maximum allowed size")
	}

	compressed, err := c.compress(data)
	if err != nil {
		return nil, err
	}

	if len(compressed) > MaxMessageSize {
		return nil, errors.New("compressed message size exceeds maximum allowed size")
	}

	encrypted, err := c.encrypt(compressed)
	if err != nil {
		return nil, err
	}

	if len(encrypted) > MaxMessageSize {
		return nil, errors.New("encrypted message size exceeds maximum allowed size")
	}

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

	log.Info().Msgf("published message with ID %s", id)
	return nil
}

func (c *Conf) serialize(message []byte) ([]byte, error) {
	if !c.Serialize {
		return message, nil
	}

	return json.Marshal(message)
}

func (c *Conf) encrypt(data []byte) ([]byte, error) {
	if c.encryptor == nil {
		return data, nil
	}

	return c.encryptor.Encrypt(data)
}

func (c *Conf) compress(data []byte) ([]byte, error) {
	if !c.Compress {
		return data, nil
	}

	compressed, err := zstd.CompressBuffer(data)
	return compressed, err
}
