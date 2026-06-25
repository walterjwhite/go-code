package transport

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, event any) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, subscription string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(ctx context.Context, message []byte) error
