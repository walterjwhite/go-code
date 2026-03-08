package icims

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

const (
	maxEmailLength    = 254
	maxPasswordLength = 128
	maxQAKeyLength    = 500
	maxQALValueLength = 5000
	rateLimitDelay    = 500 * time.Millisecond
)

type IcimsATS struct {
	mu           sync.Mutex
	lastActionAt time.Time
}

func (i *IcimsATS) GetName() string {
	return "icims"
}

func (i *IcimsATS) applyRateLimit() {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.lastActionAt.IsZero() {
		elapsed := time.Since(i.lastActionAt)
		if elapsed < rateLimitDelay {
			time.Sleep(rateLimitDelay - elapsed)
		}
	}
	i.lastActionAt = time.Now()
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if len(email) > maxEmailLength {
		return fmt.Errorf("email exceeds maximum length of %d", maxEmailLength)
	}
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, email)
	if err != nil {
		return fmt.Errorf("email validation error: %w", err)
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) > maxPasswordLength {
		return fmt.Errorf("password exceeds maximum length of %d", maxPasswordLength)
	}
	return nil
}

func validateQAMap(qaMap map[string]string) error {
	if qaMap == nil {
		return errors.New("qaMap cannot be nil")
	}
	for key, value := range qaMap {
		if len(key) > maxQAKeyLength {
			return fmt.Errorf("qaMap key exceeds maximum length of %d", maxQAKeyLength)
		}
		if len(value) > maxQALValueLength {
			return fmt.Errorf("qaMap value exceeds maximum length of %d", maxQALValueLength)
		}
		if strings.Contains(key, "<script") || strings.Contains(value, "<script") {
			return errors.New("potential XSS injection detected in qaMap")
		}
		if strings.Contains(key, "--") || strings.Contains(value, "--") {
			return errors.New("potential SQL injection detected in qaMap")
		}
	}
	return nil
}

func sanitizePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path cannot be empty")
	}
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return "", errors.New("path traversal detected: relative paths are not allowed")
	}
	if !filepath.IsAbs(cleanPath) {
		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		cleanPath = absPath
	}
	return cleanPath, nil
}

func (i *IcimsATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	if executor == nil {
		return errors.New("executor cannot be nil")
	}
	if account == nil {
		return errors.New("account cannot be nil")
	}
	if account.Email != "" {
		if err := validateEmail(account.Email); err != nil {
			return fmt.Errorf("invalid account email: %w", err)
		}
	}
	return nil
}

func (i *IcimsATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	if executor == nil {
		return errors.New("executor cannot be nil")
	}

	i.applyRateLimit()

	if err := validateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

func (i *IcimsATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	if executor == nil {
		return errors.New("executor cannot be nil")
	}
	sanitizedPath, err := sanitizePath(resumePath)
	if err != nil {
		return fmt.Errorf("invalid resume path: %w", err)
	}
	if err := validateQAMap(qaMap); err != nil {
		return fmt.Errorf("invalid qaMap: %w", err)
	}
	_ = sanitizedPath // Use sanitized path in actual implementation
	return nil
}
