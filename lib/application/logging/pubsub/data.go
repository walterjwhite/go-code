package pubsub

import (
	"context"
	"github.com/walterjwhite/go-code/lib/net/google"
)

type Publisher interface {
	Publish(topic string, message []byte) error
	Init(pctx context.Context)
	Cancel()
}

type PubsubWriter struct {
	Conf      *google.Conf
	TopicName string
	Level     string
}

func (w *PubsubWriter) Write(p []byte) (n int, err error) {
	return len(p), w.Conf.Publish(w.TopicName, p)
}

func (w *PubsubWriter) Close() {
	w.Conf.Cancel()
}
