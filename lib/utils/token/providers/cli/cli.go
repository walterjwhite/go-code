package cli

import (
	"errors"
	"flag"
	"regexp"
	"sync"
	"unicode"
)

var tokenPattern = regexp.MustCompile(`^[0-9]{6}$`)

type Provider struct {
	token string
	mu    sync.RWMutex
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("token-provider", flag.ContinueOnError)
	tokenFlag := fs.String("t", "", "RSA Token")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := p.validateToken(*tokenFlag); err != nil {
		return err
	}

	p.mu.Lock()
	p.token = *tokenFlag
	p.mu.Unlock()

	return nil
}

func (p *Provider) validateToken(token string) error {
	if len(token) == 0 {
		return ErrEmptyToken
	}
	if !tokenPattern.MatchString(token) {
		return ErrInvalidTokenFormat
	}
	hasDigit := false
	for _, r := range token {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return ErrInvalidTokenFormat
	}
	return nil
}

func (p *Provider) Get() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.token
}

var (
	ErrEmptyToken         = errors.New("token cannot be empty")
	ErrInvalidTokenFormat = errors.New("token must be 8-12 alphanumeric characters with at least one digit")
)


