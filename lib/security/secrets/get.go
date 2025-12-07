package secrets

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func Get(secretName string) string {
	if strings.TrimSpace(secretName) == "" {
		log.Warn().Msg("secrets.Get called with empty secretName")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "secrets", "get", "-out=stdout", secretName)
	cmd.Env = append(os.Environ(), "_FORCE_INTERACTIVE=1")

	out, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msgf("secrets.Get failed for key")
		return ""
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		log.Warn().Msgf("secrets.Get returned empty value for %s", secretName)
	}

	return s
}
