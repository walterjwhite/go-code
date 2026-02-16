package pubsub

import (
	"context"
)

type Publisher interface {
	Publish(topic string, message []byte) error
	Init(pctx context.Context)
	Cancel()
}

type PubsubWriter struct {
	Publisher Publisher `yaml:"-"`
	TopicName string    `yaml:"TopicName"`
	Level     string    `yaml:"Level"`
}

func (w *PubsubWriter) Init(ctx context.Context, publisher Publisher) {
	w.Publisher = publisher
	if w.Publisher != nil {
		w.Publisher.Init(ctx)
	}
}

func (w *PubsubWriter) Write(p []byte) (n int, err error) {
	if w.Publisher == nil {
		return len(p), nil
	}
	return len(p), w.Publisher.Publish(w.TopicName, p)
}

func (w *PubsubWriter) Close() {
	if w.Publisher != nil {
		w.Publisher.Cancel()
	}
}
