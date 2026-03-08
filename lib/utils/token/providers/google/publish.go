package google

import (
	"errors"
	"fmt"
	"regexp"
)

var tokenPattern = regexp.MustCompile(`^[0-9]{6}$`)

var (
	ErrInvalidTokenFormat = errors.New("token must contain only alphanumeric characters and be at least 8 characters long")
)

func validateToken(token string) error {
	if !tokenPattern.MatchString(token) {
		return ErrInvalidTokenFormat
	}
	return nil
}

func (p *Provider) PublishToken(token string) error {
	if err := validateToken(token); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	return p.Conf.Publish(p.TokenTopicName, []byte(token))
}

func (p *Provider) PublishStatus(status string, successful bool) error {
	return p.Conf.Publish(p.StatusTopicName, fmt.Appendf(nil, "%s|%v", status, successful))
}
