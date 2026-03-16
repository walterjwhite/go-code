package run

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"os/exec"
	"slices"
	"strings"
	"time"
)

type Exec struct {
	Command   string
	Arguments []string
}

var allowedCommands = map[string]bool{
	"sleep":  true,
	"test":   true,
	"true":   true,
	"false":  true,
	"pwd":    true,
	"whoami": true,
	"date":   true,
}

func (e *Exec) Do(context.Context) error {
	if !isCommandAllowed(e.Command) {
		return errors.New("command not allowlisted: " + e.Command)
	}

	if slices.ContainsFunc(e.Arguments, containsDangerousPatterns) {
		return errors.New("argument contains dangerous patterns: potential command injection detected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.Command, e.Arguments...) // #nosec G204 - command is validated by validateCommand

	log.Info().Msgf("running: %s %v", e.Command, e.Arguments)
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msgf("command failed: %s", e.Command)
		return err
	}

	log.Info().Msgf("done running: %s %v", e.Command, e.Arguments)
	return nil
}

func isCommandAllowed(cmd string) bool {
	return allowedCommands[cmd]
}

func containsDangerousPatterns(arg string) bool {
	dangerousPatterns := []string{";", "|", "&", "$", "`", "(", ")", "<", ">", "\n", "\r"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(arg, pattern) {
			return true
		}
	}
	return false
}
