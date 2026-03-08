package secrets

import (
	"context"

	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func Get(secretName string) string {
	if strings.TrimSpace(secretName) == "" {
		log.Warn().Msg("secrets.Get called with empty secretName")
		return ""
	}

	if !isValidSecretName(secretName) {
		log.Warn().Msg("secrets.Get called with invalid secret name format")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "secrets", "-out=stdout", "get", secretName)

	out, err := cmd.Output()
	if err != nil {
		log.Error().Msg("failed to retrieve secret from secrets service")
		return ""
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		log.Warn().Msg("secrets service returned empty value")
	}

	return s
}

func isValidSecretName(secretName string) bool {
	if len(secretName) == 0 || len(secretName) > 256 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`)
	if !re.MatchString(secretName) {
		return false
	}

	if strings.HasPrefix(secretName, "/") || strings.HasSuffix(secretName, "/") {
		return false
	}

	if strings.Contains(secretName, "//") {
		return false
	}

	if strings.Contains(secretName, "..") {
		return false
	}

	return true
}
