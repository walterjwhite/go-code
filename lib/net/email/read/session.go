package write

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
)

type EmailSession struct {
	EmailAccount *email.EmailAccount

	client *client.Client
}

func sanitizeConnectionError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	if containsSensitiveInfo(errStr) {
		return errors.New("failed to establish secure connection")
	}

	return err
}

func containsSensitiveInfo(errStr string) bool {
	sensitivePatterns := []string{
		"password", "credential", "secret", "token",
		"private", "key", "certificate",
	}

	for _, pattern := range sensitivePatterns {
		if containsIgnoreCase(errStr, pattern) {
			return true
		}
	}

	return false
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && containsLower(s, substr)
}

func containsLower(s, substr string) bool {
	s = toLowerASCII(s)
	substr = toLowerASCII(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func sanitizeLoginError(err error) error {
	if err == nil {
		return nil
	}

	return errors.New("authentication failed")
}

func New(a *email.EmailAccount) (*EmailSession, error) {
	if a == nil {
		return nil, errors.New("email account cannot be nil")
	}

	if a.ImapServer == nil {
		return nil, errors.New("IMAP server configuration cannot be nil")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         a.ImapServer.Host,
		MinVersion:         tls.VersionTLS12,
	}

	c, err := client.DialTLS(fmt.Sprintf("%v:%v", a.ImapServer.Host, a.ImapServer.Port), tlsConfig)
	if err != nil {
		return nil, sanitizeConnectionError(err)
	}

	err = c.Login(a.Username, a.Password)
	if err != nil {
		_ = c.Logout()
		return nil, sanitizeLoginError(err)
	}

	log.Info().Msg("successfully logged in")

	return &EmailSession{EmailAccount: a, client: c}, nil
}

func (s *EmailSession) Close() {
	if s.client == nil {
		return
	}

	if err := s.client.Logout(); err != nil {
		log.Warn().Msg("failed to logout from email session")
	}
	s.client = nil
}
