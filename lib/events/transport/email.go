package transport

import (
	"context"
	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	"github.com/walterjwhite/go-code/lib/net/messaging"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"gopkg.in/gomail.v2"
)

type EmailPublisher struct {
	dialer     *gomail.Dialer
	from       string
	processor  *messaging.MessageProcessor
	serializer serialization.Serializer
}

func NewEmailPublisher(
	host string,
	port int,
	username string,
	password string,
	from string,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor any,
	enableCompression bool,
	enableEncryption bool,
) *EmailPublisher {
	dialer := gomail.NewDialer(host, port, username, password)

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

	return &EmailPublisher{
		dialer:     dialer,
		from:       from,
		processor:  processor,
		serializer: serializer,
	}
}

func (p *EmailPublisher) Publish(ctx context.Context, topic string, event any) error {
	if p.dialer == nil {
		return fmt.Errorf("email dialer not initialized")
	}

	serialized, err := p.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	processed, err := p.processor.Process(serialized)
	if err != nil {
		return fmt.Errorf("failed to process event: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", p.from)
	m.SetHeader("To", topic)
	m.SetHeader("Subject", "Event Notification")
	m.SetBody("text/plain", string(processed))

	if err := p.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", topic, err)
	}

	return nil
}

func (p *EmailPublisher) Close() error {
	return nil
}

type EmailSubscriber struct {
	host       string
	username   string
	password   string
	processor  *messaging.MessageProcessor
	serializer serialization.Serializer
}

func NewEmailSubscriber(
	host string,
	username string,
	password string,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor any,
	enableCompression bool,
	enableEncryption bool,
) *EmailSubscriber {
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

	return &EmailSubscriber{
		host:       host,
		username:   username,
		password:   password,
		processor:  processor,
		serializer: serializer,
	}
}

func (s *EmailSubscriber) Subscribe(ctx context.Context, subscription string, handler MessageHandler) error {
	if subscription == "" {
		subscription = "INBOX"
	}

	if strings.ContainsAny(subscription, "/\\") {
		return fmt.Errorf("invalid subscription folder: %s", subscription)
	}

	return fmt.Errorf("email subscription implementation requires IMAP client setup")
}

func (s *EmailSubscriber) Close() error {
	return nil
}
