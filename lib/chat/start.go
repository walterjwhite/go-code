package chat

import "context"

func (c *Client) Start(ctx context.Context) {
	c.startOnce.Do(func() {
		go c.startSubscription(ctx)
	})
}

func (c *Client) startSubscription(ctx context.Context) {
	if ctx == nil {
		c.pushError(ErrNilContext)
		return
	}
	c.pubsub.Subscribe(
		c.config.ResponseSubscription,
		c.newSubscriber(ctx),
	)
}

func (c *Client) newSubscriber(ctx context.Context) *responseSubscriber {
	return &responseSubscriber{
		ctx:        ctx,
		serializer: c.serializer,
		responses:  c.responses,
		errors:     c.errors,
	}
}

func (c *Client) Responses() <-chan ServerMessage { return c.responses }

func (c *Client) Errors() <-chan error { return c.errors }

func (c *Client) pushError(err error) {
	select {
	case c.errors <- err:
	default:
	}
}
