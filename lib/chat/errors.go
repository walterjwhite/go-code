package chat

import "errors"

var (
	ErrNilPubSub              = errors.New("pubsub is required")
	ErrNilSerializer          = errors.New("serializer is required")
	ErrNilContext             = errors.New("context is required")
	ErrPublishTopicRequired   = errors.New("publish topic is required")
	ErrResponseTopicRequired  = errors.New("response topic is required")
	ErrResponseSubRequired    = errors.New("response subscription is required")
	ErrEmptyIP                = errors.New("ip is required")
	ErrEmptyMessage           = errors.New("message is required")
	ErrRateLimitWindowInvalid = errors.New("rate limit window must be greater than zero")
	ErrRateLimitCountInvalid  = errors.New("rate limit count must be greater than zero")
	ErrRateLimited            = errors.New("rate limit exceeded for ip")
)
