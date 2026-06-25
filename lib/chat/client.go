package chat

import (
	"context"

	"github.com/walterjwhite/go-code/lib/io/serialization"
)

func NewClient(config Config, pubsub PubSub, serializer serialization.Serializer) (*Client, error) {
	if err := validateDependencies(pubsub, serializer); err != nil {
		return nil, err
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return newClient(config, pubsub, serializer), nil
}

func validateDependencies(pubsub PubSub, serializer serialization.Serializer) error {
	if pubsub == nil {
		return ErrNilPubSub
	}
	if serializer == nil {
		return ErrNilSerializer
	}
	return nil
}

func newClient(config Config, pubsub PubSub, serializer serialization.Serializer) *Client {
	size := responseBufferSize(config)
	limiter := NewIPRateLimiter(config.RateLimit)
	limiter.Start(context.Background(), config.RateLimit.Window)
	return &Client{
		config:     config,
		pubsub:     pubsub,
		serializer: serializer,
		limiter:    limiter,
		responses:  make(chan ServerMessage, size),
		errors:     make(chan error, size),
	}
}
