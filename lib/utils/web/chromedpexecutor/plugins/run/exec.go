package run

import (
	"context"
	"github.com/rs/zerolog/log"
	"os/exec"
	"time"
)

type Exec struct {
	Command   string
	Arguments []string
}

func (e *Exec) Do(context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.Command, e.Arguments...)

	log.Info().Msgf("running: %s %v", e.Command, e.Arguments)
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msgf("command failed: %s", e.Command)
		return err
	}

	log.Info().Msgf("done running: %s %v", e.Command, e.Arguments)
	return nil
}
