package transport

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	"github.com/walterjwhite/go-code/lib/net/messaging"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"gopkg.in/gomail.v2"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type EmailPublisher struct {
	dialer     *gomail.Dialer
	from       string
	processor  *messaging.MessageProcessor
	serializer serialization.Serializer
}

type EmailPublisherConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	From              string
	Serializer        serialization.Serializer
	Compressor        compression.Compressor
	Encryptor         encryption.Encryptor
	EnableCompression bool
	EnableEncryption  bool
	UseTLS            bool
}

func NewEmailPublisher(
	host string,
	port int,
	username string,
	password string,
	from string,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor encryption.Encryptor,
	enableCompression bool,
	enableEncryption bool,
) *EmailPublisher {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.TLSConfig = nil // Uses default TLS config

	processor := messaging.NewMessageProcessor(
		serializer,
		compressor,
		encryptor,
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

func NewEmailPublisherWithConfig(cfg EmailPublisherConfig) *EmailPublisher {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	if cfg.UseTLS {
		dialer.TLSConfig = nil // Uses default TLS config
	}

	processor := messaging.NewMessageProcessor(
		cfg.Serializer,
		cfg.Compressor,
		cfg.Encryptor,
		false, // serialization handled separately
		cfg.EnableCompression,
		cfg.EnableEncryption,
	)

	return &EmailPublisher{
		dialer:     dialer,
		from:       cfg.From,
		processor:  processor,
		serializer: cfg.Serializer,
	}
}

func validateEmailAddress(email string) error {
	if email == "" {
		return fmt.Errorf("email address cannot be empty")
	}

	if strings.ContainsAny(email, "\r\n") {
		return fmt.Errorf("invalid email address: contains control characters")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email address format: %s", email)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

func (p *EmailPublisher) Publish(ctx context.Context, topic string, event any) error {
	if p.dialer == nil {
		return fmt.Errorf("email dialer not initialized")
	}

	if err := validateEmailAddress(topic); err != nil {
		return fmt.Errorf("invalid recipient: %w", err)
	}

	if err := validateEmailAddress(p.from); err != nil {
		return fmt.Errorf("invalid sender: %w", err)
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
		return fmt.Errorf("failed to send email: %w", err)
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

type EmailSubscriberConfig struct {
	Host              string
	Username          string
	Password          string
	Serializer        serialization.Serializer
	Compressor        compression.Compressor
	Encryptor         encryption.Encryptor
	EnableCompression bool
	EnableEncryption  bool
	UseTLS            bool
}

func NewEmailSubscriber(
	host string,
	username string,
	password string,
	serializer serialization.Serializer,
	compressor compression.Compressor,
	encryptor encryption.Encryptor,
	enableCompression bool,
	enableEncryption bool,
) *EmailSubscriber {
	processor := messaging.NewMessageProcessor(
		serializer,
		compressor,
		encryptor,
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

func NewEmailSubscriberWithConfig(cfg EmailSubscriberConfig) *EmailSubscriber {
	processor := messaging.NewMessageProcessor(
		cfg.Serializer,
		cfg.Compressor,
		cfg.Encryptor,
		false, // serialization handled separately
		cfg.EnableCompression,
		cfg.EnableEncryption,
	)

	return &EmailSubscriber{
		host:       cfg.Host,
		username:   cfg.Username,
		password:   cfg.Password,
		processor:  processor,
		serializer: cfg.Serializer,
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
