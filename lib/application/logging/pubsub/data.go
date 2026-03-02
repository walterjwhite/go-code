package pubsub

import (
	"context"
	"fmt"
	"regexp"
	"sync"
)

type Publisher interface {
	Publish(topic string, message []byte) error
	Init(pctx context.Context) error
	Cancel()
}

type PubsubWriter struct {
	mu        sync.RWMutex
	Publisher Publisher `yaml:"-"`
	TopicName string    `yaml:"TopicName"`
	Level     string    `yaml:"Level"`
}

var topicNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,255}$`)

func validateTopicName(topic string) error {
	if topic == "" {
		return fmt.Errorf("topic name cannot be empty")
	}
	if !topicNameRegex.MatchString(topic) {
		return fmt.Errorf("invalid topic name: must contain only alphanumeric characters, dashes, and underscores, and be between 1-255 characters")
	}
	return nil
}

func (w *PubsubWriter) Init(ctx context.Context, publisher Publisher) error {
	if err := validateTopicName(w.TopicName); err != nil {
		return fmt.Errorf("pubsub writer initialization failed: %w", err)
	}

	w.mu.Lock()
	w.Publisher = publisher
	w.mu.Unlock()

	if publisher != nil {
		return publisher.Init(ctx)
	}
	return nil
}

func (w *PubsubWriter) Write(p []byte) (n int, err error) {
	w.mu.RLock()
	publisher := w.Publisher
	topicName := w.TopicName
	w.mu.RUnlock()

	if publisher == nil {
		return 0, fmt.Errorf("publisher not initialized: cannot write to topic %q", topicName)
	}
	err = publisher.Publish(topicName, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *PubsubWriter) Close() {
	w.mu.RLock()
	publisher := w.Publisher
	w.mu.RUnlock()

	if publisher != nil {
		publisher.Cancel()
	}
}
