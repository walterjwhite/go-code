package run

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os/exec"
)

type Exec struct {
	Command   string
	Arguments []string
}

func (e *Exec) Do(context.Context) error {
	cmd := exec.Command(e.Command, e.Arguments...)

	log.Info().Msgf("running: %s %s", e.Command, e.Arguments)
	logging.Panic(cmd.Run())

	log.Info().Msgf("done running: %s %s", e.Command, e.Arguments)
	return nil
}
