package chat

import (
	"context"
	"sync"
	"time"

	"github.com/walterjwhite/go-code/lib/io/serialization"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type PubSub interface {
	Publish(topic string, message []byte) error
	Subscribe(subscriptionName string, ms google_pubsub.MessageSubscriber)
}

type Config struct {
	PublishTopic         string
	ResponseTopic        string
	ResponseSubscription string
	ResponseBuffer       int
	RateLimit            RateLimitConfig
}

type RateLimitConfig struct {
	Messages int
	Window   time.Duration
	MaxIPs int
}

const DefaultMaxTrackedIPs = 10_000

type ClientMessage struct {
	IP            string            `json:"ip"`
	Message       string            `json:"message"`
	Channel       string            `json:"channel,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SentAt        time.Time         `json:"sent_at"`
}

type ServerMessage struct {
	Channel       string            `json:"channel,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Message       string            `json:"message"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SentAt        time.Time         `json:"sent_at"`
}

type Client struct {
	config     Config
	pubsub     PubSub
	serializer serialization.Serializer
	limiter    *IPRateLimiter
	responses  chan ServerMessage
	errors     chan error
	startOnce  sync.Once
}

type responseSubscriber struct {
	ctx        context.Context
	serializer serialization.Serializer
	responses  chan<- ServerMessage
	errors     chan<- error
}

type IPRateLimiter struct {
	config  RateLimitConfig
	maxIPs  int
	mu      sync.Mutex
	windows map[string]ipWindow
}

type ipWindow struct {
	start time.Time
	count int
}
