package pubsub

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
)

type Publisher interface {
	Publish(topic string, message []byte) error
	Init(pctx context.Context) error
	Cancel()
}

type PubSubWriter struct {
	mu        sync.RWMutex
	Publisher Publisher `yaml:"-"`
	TopicName string    `yaml:"TopicName"`
	Level     string    `yaml:"Level"`
}

var topicNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,255}$`)

const maxMessageSize = 10 * 1024 * 1024

func validateTopicName(topic string) error {
	if topic == "" {
		return fmt.Errorf("topic name cannot be empty")
	}
	if !topicNameRegex.MatchString(topic) {
		return fmt.Errorf("invalid topic name: must contain only alphanumeric characters, dashes, and underscores, and be between 1-255 characters")
	}
	return nil
}

func validateMessage(message []byte) error {
	if len(message) == 0 {
		return errors.New("message cannot be empty")
	}
	if len(message) > maxMessageSize {
		return fmt.Errorf("message size %d exceeds maximum allowed size of %d bytes", len(message), maxMessageSize)
	}
	return nil
}

func (w *PubSubWriter) getPublisherAndTopic() (Publisher, string) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Publisher, w.TopicName
}

func (w *PubSubWriter) setPublisher(publisher Publisher) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Publisher = publisher
}

func (w *PubSubWriter) getPublisher() Publisher {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Publisher
}

func (w *PubSubWriter) Init(ctx context.Context, publisher Publisher) error {
	if err := validateTopicName(w.TopicName); err != nil {
		return fmt.Errorf("PubSub writer initialization failed: %w", err)
	}

	w.setPublisher(publisher)

	if publisher != nil {
		return publisher.Init(ctx)
	}
	return nil
}

func (w *PubSubWriter) Write(p []byte) (n int, err error) {
	if err := validateMessage(p); err != nil {
		return 0, fmt.Errorf("message validation failed: %w", err)
	}

	publisher, topicName := w.getPublisherAndTopic()

	if publisher == nil {
		return 0, fmt.Errorf("publisher not initialized: cannot write to topic %q", topicName)
	}

	err = publisher.Publish(topicName, p)
	if err != nil {
		return 0, fmt.Errorf("failed to publish message to topic %q: %w", topicName, err)
	}

	return len(p), nil
}

func (w *PubSubWriter) Close() {
	publisher := w.getPublisher()
	if publisher != nil {
		publisher.Cancel()
	}
}
