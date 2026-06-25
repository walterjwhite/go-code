package chat

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) Send(ctx context.Context, message ClientMessage) error {
	if err := validateSend(ctx, message); err != nil {
		return err
	}
	if !c.limiter.Allow(message.IP) {
		return ErrRateLimited
	}
	return c.publish(message.withTimestamp())
}

func validateSend(ctx context.Context, message ClientMessage) error {
	if ctx == nil {
		return ErrNilContext
	}
	if message.IP == "" {
		return ErrEmptyIP
	}
	if message.Message == "" {
		return ErrEmptyMessage
	}
	return nil
}

func (c *Client) publish(message ClientMessage) error {
	payload, err := c.serializer.Serialize(message)
	if err != nil {
		return fmt.Errorf("serialize chat message: %w", err)
	}
	if err := c.pubsub.Publish(c.config.PublishTopic, payload); err != nil {
		return fmt.Errorf("publish chat message: %w", err)
	}
	return nil
}

func (m ClientMessage) withTimestamp() ClientMessage {
	if m.SentAt.IsZero() {
		m.SentAt = time.Now().UTC()
	}
	return m
}
