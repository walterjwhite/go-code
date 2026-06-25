package google

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/walterjwhite/go-code/lib/application"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

const (
	DefaultMaxMessageSize = 10 * 1024 * 1024 // 10MB - GCP Pub/Sub limit
	DefaultRateLimit      = 100              // messages per second
	RateLimitWindow       = time.Second
)

var (
	ErrEmptyMessage         = errors.New("message cannot be empty")
	ErrMessageTooLarge      = errors.New("message exceeds maximum allowed size")
	ErrInvalidMessage       = errors.New("message contains invalid characters")
	ErrNilContext           = errors.New("context cannot be nil")
	ErrContextCancelled     = errors.New("context cancelled")
	ErrContextDeadline      = errors.New("context deadline exceeded")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrNotInitialized       = errors.New("provider not initialized")
	ErrMissingTopicName     = errors.New("topic name not configured")
	ErrInvalidTopicName     = errors.New("topic name contains invalid characters")
	ErrInvalidSubscription  = errors.New("subscription name contains invalid characters")
	safeMessagePattern      = regexp.MustCompile(`^[\x20-\x7E\x09\x0A\x0D]*$`)     // Printable ASCII + whitespace
	safeResourceNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{0,254}$`) // Valid resource name pattern
)

type Provider struct {
	TopicName        string
	SubscriptionName string
	MaxMessageSize   int
	RateLimit        int

	Conf *google_pubsub.Conf

	mu       sync.Mutex
	lastCall time.Time
	count    int
}

func New(ctx context.Context) (*Provider, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	provider := &Provider{
		MaxMessageSize: DefaultMaxMessageSize,
		RateLimit:      DefaultRateLimit,
	}
	application.Load(provider)

	if provider.Conf == nil {
		return nil, ErrNotInitialized
	}

	if err := provider.Conf.Init(ctx); err != nil {
		return nil, ErrNotInitialized
	}

	return provider, nil
}

func (p *Provider) String() string {
	return "Provider: {configured}"
}

func (p *Provider) Publish(ctx context.Context, message []byte) error {
	if ctx == nil {
		return ErrNilContext
	}

	if err := checkContext(ctx); err != nil {
		return err
	}

	if p.Conf == nil {
		return ErrNotInitialized
	}

	if err := validateResourceNames(p.TopicName, p.SubscriptionName); err != nil {
		return err
	}

	if err := p.checkRateLimit(); err != nil {
		return err
	}

	if err := validateMessage(message, p.MaxMessageSize); err != nil {
		return err
	}

	return p.Conf.Publish(p.TopicName, message)
}

func (p *Provider) checkRateLimit() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if now.Sub(p.lastCall) > RateLimitWindow {
		p.count = 0
		p.lastCall = now
	}

	if p.count >= p.RateLimit {
		return ErrRateLimitExceeded
	}

	p.count++
	return nil
}

func validateMessage(message []byte, maxSize int) error {
	if len(message) == 0 {
		return ErrEmptyMessage
	}

	if len(message) > maxSize {
		return ErrMessageTooLarge
	}

	if !utf8.Valid(message) {
		return ErrInvalidMessage
	}

	if !safeMessagePattern.Match(message) {
		return ErrInvalidMessage
	}

	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return ErrContextDeadline
		}
		return ErrContextCancelled
	default:
		return nil
	}
}

func validateResourceNames(topicName, subscriptionName string) error {
	if topicName == "" {
		return ErrMissingTopicName
	}

	if !safeResourceNamePattern.MatchString(topicName) {
		return ErrInvalidTopicName
	}

	if subscriptionName != "" && !safeResourceNamePattern.MatchString(subscriptionName) {
		return ErrInvalidSubscription
	}

	return nil
}
