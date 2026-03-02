package transport

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	"github.com/walterjwhite/go-code/lib/net/messaging"
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

type GooglePubSubPublisher struct {
	client     *pubsub.Client
	processor  *messaging.MessageProcessor
	serializer serialization.Serializer
}

func NewGooglePubSubPublisher(
	client *pubsub.Client,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor any,
	enableCompression bool,
	enableEncryption bool,
) *GooglePubSubPublisher {
	var encInt encryption.Encryptor
	if encryptor != nil {
		var ok bool
		encInt, ok = encryptor.(encryption.Encryptor)
		if !ok {
			panic("encryptor must implement encryption.Encryptor interface")
		}
	}

	processor := messaging.NewMessageProcessor(
		serializer,
		compressor,
		encInt,
		false, // serialization handled separately
		enableCompression,
		enableEncryption,
	)

	return &GooglePubSubPublisher{
		client:     client,
		processor:  processor,
		serializer: serializer,
	}
}

func (p *GooglePubSubPublisher) Publish(ctx context.Context, topic string, event any) error {
	if p.client == nil {
		return fmt.Errorf("publisher not initialized")
	}

	serialized, err := p.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	processed, err := p.processor.Process(serialized)
	if err != nil {
		return fmt.Errorf("failed to process event: %w", err)
	}

	publisher := p.client.Publisher(topic)
	defer publisher.Stop()
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: processed,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, err)
	}

	return nil
}

func (p *GooglePubSubPublisher) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

type GooglePubSubSubscriber struct {
	client     *pubsub.Client
	processor  *messaging.MessageProcessor
	serializer serialization.Serializer
}

func NewGooglePubSubSubscriber(
	client *pubsub.Client,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor any,
	enableCompression bool,
	enableEncryption bool,
) *GooglePubSubSubscriber {
	var encInt encryption.Encryptor
	if encryptor != nil {
		var ok bool
		encInt, ok = encryptor.(encryption.Encryptor)
		if !ok {
			panic("encryptor must implement encryption.Encryptor interface")
		}
	}

	processor := messaging.NewMessageProcessor(
		serializer,
		compressor,
		encInt,
		false, // serialization handled separately
		enableCompression,
		enableEncryption,
	)

	return &GooglePubSubSubscriber{
		client:     client,
		processor:  processor,
		serializer: serializer,
	}
}

func (s *GooglePubSubSubscriber) Subscribe(ctx context.Context, subscription string, handler MessageHandler) error {
	if s.client == nil {
		return fmt.Errorf("subscriber not initialized")
	}

	sub := s.client.Subscriber(subscription)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		unprocessed, err := s.processor.Unprocess(msg.Data)
		if err != nil {
			msg.Nack()
			return
		}

		if err := handler(ctx, unprocessed); err != nil {
			msg.Nack()
			return
		}

		msg.Ack()
	})

	return err
}

func (s *GooglePubSubSubscriber) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
