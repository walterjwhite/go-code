package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub/v2"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	googleconf "github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"google.golang.org/api/option"
)

type MessageMetadata struct {
	ID      string
	Profile string
	index   int
}

type PubSubRuntime struct {
	ctx       context.Context
	client    *pubsub.Client
	encryptor encryption.Encryptor
	compress  bool
	serialize bool
}

func NewPubSubRuntime(ctx context.Context, conf *googleconf.Conf) (*PubSubRuntime, error) {
	if len(conf.CredentialsFile) == 0 {
		return nil, fmt.Errorf("credentials file path is empty: must be configured")
	}

	if _, err := os.Stat(conf.CredentialsFile); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("credentials file not found")
		}
		return nil, fmt.Errorf("credentials file validation failed")
	}

	client, err := pubsub.NewClient(ctx, conf.ProjectId, option.WithAuthCredentialsFile(option.ServiceAccount, conf.CredentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	var encryptor encryption.Encryptor
	if len(conf.EncryptionKeyFile) > 0 {
		encryptor, err = aes.NewAESFromFile(conf.EncryptionKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize encryption: %w", err)
		}
	}

	return &PubSubRuntime{
		ctx:       ctx,
		client:    client,
		encryptor: encryptor,
		compress:  conf.Compress,
		serialize: conf.Serialize,
	}, nil
}

func (r *PubSubRuntime) Subscribe(subscriptionName string, e *Executor) {
	sub := r.client.Subscriber(subscriptionName)
	log.Info().Msgf("subscribed to: %s", subscriptionName)

	err := sub.Receive(r.ctx, func(ctx context.Context, m *pubsub.Message) {
		profile := ""
		id := m.ID
		if m.Attributes != nil {
			if val, ok := m.Attributes["profile"]; ok {
				profile = val
			}
			if val, ok := m.Attributes["response-to"]; ok {
				id = val
			}
		}

		log.Info().Msgf("received message: %s", m.ID)

		metadata := &MessageMetadata{ID: id, Profile: profile}
		e.publishEvent(metadata, map[string]string{"status": "received"})

		deserialized, err := r.processIncomingMessage(m.Data, metadata, e)
		if err != nil {
			logging.Warn(err, "pubsub-exec.processIncomingMessage")
			m.Ack()
			return
		}

		e.MessageDeserializedWithMetadata(deserialized, metadata)
		m.Ack()
	})
	logging.Warn(err, "pubsub-exec.subscribe.Receive")
}

func (r *PubSubRuntime) processIncomingMessage(data []byte, metadata *MessageMetadata, e *Executor) ([]byte, error) {
	if len(data) > googleconf.MaxMessageSize {
		err := fmt.Errorf("received message size %d exceeds maximum allowed size %d", len(data), googleconf.MaxMessageSize)
		e.publishEvent(metadata, map[string]string{"status": "processing failed", "done": "true", "exitCode": "1", "output": err.Error()})
		return nil, err
	}

	decrypted, err := r.decrypt(data)
	if err != nil {
		e.publishEvent(metadata, map[string]string{"status": "processing failed", "done": "true", "exitCode": "2", "output": err.Error()})
		return nil, err
	}
	e.publishEvent(metadata, map[string]string{"status": "decrypted"})

	if len(decrypted) > googleconf.MaxMessageSize {
		err := errors.New("decrypted message size exceeds maximum allowed size")
		e.publishEvent(metadata, map[string]string{"status": "processing failed", "done": "true", "exitCode": "3", "output": err.Error()})
		return nil, err
	}

	decompressed, err := r.decompress(decrypted)
	if err != nil {
		e.publishEvent(metadata, map[string]string{"status": "processing failed", "done": "true", "exitCode": "4", "output": err.Error()})
		return nil, err
	}
	e.publishEvent(metadata, map[string]string{"status": "decompressed"})

	if len(decompressed) > googleconf.MaxMessageSize {
		err := errors.New("decompressed message size exceeds maximum allowed size")
		e.publishEvent(metadata, map[string]string{"status": "processing failed", "done": "true", "exitCode": "5", "output": err.Error()})
		return nil, err
	}

	deserialized, err := r.deserialize(decompressed)
	if err != nil {
		e.publishEvent(metadata, map[string]string{"status": "parse failed", "done": "true", "exitCode": "6", "output": err.Error()})
		return nil, err
	}

	return deserialized, nil
}

func (r *PubSubRuntime) Publish(topicName string, message []byte, attributes map[string]string) error {
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	default:
	}

	prepared, err := r.prepareMessage(message)
	if err != nil {
		return err
	}

	publisher := r.client.Publisher(topicName)
	defer publisher.Stop()

	result := publisher.Publish(r.ctx, &pubsub.Message{
		Data:       prepared,
		Attributes: attributes,
	})
	_, err = result.Get(r.ctx)
	return err
}

func (r *PubSubRuntime) prepareMessage(message []byte) ([]byte, error) {
	if len(message) > googleconf.MaxMessageSize {
		return nil, fmt.Errorf("message size %d exceeds maximum allowed size %d", len(message), googleconf.MaxMessageSize)
	}

	serialized, err := r.serializeMessage(message)
	if err != nil {
		return nil, err
	}
	if len(serialized) > googleconf.MaxMessageSize {
		return nil, errors.New("serialized message size exceeds maximum allowed size")
	}

	compressed, err := r.compressMessage(serialized)
	if err != nil {
		return nil, err
	}
	if len(compressed) > googleconf.MaxMessageSize {
		return nil, errors.New("compressed message size exceeds maximum allowed size")
	}

	encrypted, err := r.encrypt(compressed)
	if err != nil {
		return nil, err
	}
	if len(encrypted) > googleconf.MaxMessageSize {
		return nil, errors.New("encrypted message size exceeds maximum allowed size")
	}

	return encrypted, nil
}

func (r *PubSubRuntime) serializeMessage(message []byte) ([]byte, error) {
	if !r.serialize {
		return message, nil
	}

	return json.Marshal(message)
}

func (r *PubSubRuntime) deserialize(data []byte) ([]byte, error) {
	if !r.serialize {
		return data, nil
	}

	var deserialized []byte
	if err := json.Unmarshal(data, &deserialized); err != nil {
		return nil, err
	}

	return deserialized, nil
}

func (r *PubSubRuntime) encrypt(data []byte) ([]byte, error) {
	if r.encryptor == nil {
		return data, nil
	}

	return r.encryptor.Encrypt(data)
}

func (r *PubSubRuntime) decrypt(data []byte) ([]byte, error) {
	if r.encryptor == nil {
		return data, nil
	}

	return r.encryptor.Decrypt(data)
}

func (r *PubSubRuntime) compressMessage(data []byte) ([]byte, error) {
	if !r.compress {
		return data, nil
	}

	return zstd.CompressBuffer(data)
}

func (r *PubSubRuntime) decompress(data []byte) ([]byte, error) {
	if !r.compress {
		return data, nil
	}

	return zstd.DecompressBuffer(data)
}
