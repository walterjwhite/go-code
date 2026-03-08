package google

import (
	"context"
	"fmt"
	"sync"
	"time"

	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

const TokenExpirationDuration = 5 * time.Minute

type Provider struct {
	Conf *google_pubsub.Conf

	TokenTopicName        string
	TokenSubscriptionName string

	StatusTopicName        string
	StatusSubscriptionName string

	token    string
	tokenSet time.Time
	hasToken bool
	mu       sync.RWMutex
}

func (p *Provider) String() string {
	return fmt.Sprintf("Provider: {TokenTopicName: %s, TokenSubscriptionName: %s, StatusTopicName: %s, StatusSubscriptionName: %s, Conf: %s}", p.TokenTopicName, p.TokenSubscriptionName, p.StatusTopicName,
		p.StatusSubscriptionName, p.Conf)
}

func (p *Provider) Init(ctx context.Context) error {
	return p.Conf.Init(ctx)
}

func (p *Provider) Cleanup() {
	p.Conf.Cancel()
}

func (p *Provider) isTokenExpiredUnsafe() bool {
	if !p.hasToken {
		return true
	}
	return time.Since(p.tokenSet) > TokenExpirationDuration
}

func (p *Provider) IsTokenExpired() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isTokenExpiredUnsafe()
}

func (p *Provider) SetToken(token string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = token
	p.tokenSet = time.Now()
	p.hasToken = true
}

func (p *Provider) GetToken() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.isTokenExpiredUnsafe() {
		return ""
	}
	return p.token
}

func (p *Provider) ClearToken() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = ""
	p.hasToken = false
	p.tokenSet = time.Time{}
}
