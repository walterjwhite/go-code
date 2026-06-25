package main

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func (e *Executor) parseShellCommand(deserialized []byte) ([]string, error) {
	messageStr := string(deserialized)
	if len(messageStr) == 0 {
		return nil, fmt.Errorf("empty message / failed to parse as shell command")
	}

	log.Info().Msgf("received string message: %s", messageStr)
	if messageStr == "" {
		return nil, nil
	}

	var args []string
	var current strings.Builder
	var inQuotes bool
	var escapeNext bool

	for _, r := range messageStr {
		switch {
		case escapeNext:
			current.WriteRune(r)
			escapeNext = false
		case r == '\\':
			escapeNext = true
		case r == '"':
			inQuotes = !inQuotes
		case r == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if inQuotes {
		log.Warn().Msgf("unmatched quotes in command: %s", messageStr)
		return nil, nil
	}

	if escapeNext {
		log.Warn().Msgf("dangling escape character in command: %s", messageStr)
		return nil, nil
	}

	return args, nil
}
