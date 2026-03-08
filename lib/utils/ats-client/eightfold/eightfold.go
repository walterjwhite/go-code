package eightfold

import (
	"fmt"
	"strings"
	"sync"
	"time"

	atsclient "github.com/walterjwhite/go-code/lib/utils/ats-client"
)

type Account = atsclient.Account

const (
	maxInputLength = 1024
	maxEmailLength = 254
	rateLimitDelay = 500 * time.Millisecond
)

type EightfoldATS struct {
	mu           sync.Mutex
	lastActionAt time.Time
}

func (e *EightfoldATS) GetName() string {
	return "eightfold"
}

func (e *EightfoldATS) applyRateLimit() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.lastActionAt.IsZero() {
		elapsed := time.Since(e.lastActionAt)
		if elapsed < rateLimitDelay {
			time.Sleep(rateLimitDelay - elapsed)
		}
	}
	e.lastActionAt = time.Now()
}

func validateInput(pattern, value string) error {
	if len(value) > maxInputLength {
		return fmt.Errorf("input exceeds maximum length of %d", maxInputLength)
	}
	if len(pattern) > maxInputLength {
		return fmt.Errorf("pattern exceeds maximum length of %d", maxInputLength)
	}

	dangerousChars := []string{"<script", "</script>", "javascript:", "data:", "vbscript:"}
	for _, char := range dangerousChars {
		if strings.Contains(strings.ToLower(value), char) || strings.Contains(strings.ToLower(pattern), char) {
			return fmt.Errorf("input contains potentially dangerous script content")
		}
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if len(email) > maxEmailLength {
		return fmt.Errorf("email exceeds maximum length of %d", maxEmailLength)
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (e *EightfoldATS) RegisterAccount(executor *atsclient.Executor, account *Account) error {
	if executor == nil {
		return fmt.Errorf("executor cannot be nil")
	}
	if account == nil {
		return fmt.Errorf("account cannot be nil")
	}

	if account.Email != "" {
		if err := validateEmail(account.Email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
	}

	fmt.Println("Registering account on Eightfold...")
	return nil
}

func (e *EightfoldATS) LoginAccount(executor *atsclient.Executor, email, password string) error {
	if executor == nil {
		return fmt.Errorf("executor cannot be nil")
	}

	e.applyRateLimit()

	if err := validateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	fmt.Println("Logging into Eightfold...")
	return nil
}

func (e *EightfoldATS) ApplyForJob(executor *atsclient.Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error {
	if executor == nil {
		return fmt.Errorf("executor cannot be nil")
	}

	for pattern, value := range qaMap {
		if err := validateInput(pattern, value); err != nil {
			return fmt.Errorf("invalid qaMap entry: %w", err)
		}
	}

	fmt.Println("Applying for job on Eightfold...")
	return nil
}
