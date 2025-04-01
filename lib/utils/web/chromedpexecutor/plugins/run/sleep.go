package run

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type Sleep struct {
	Duration time.Duration
}

func (s *Sleep) Do(context.Context) error {
	log.Info().Msgf("sleeping: %s", s.Duration)
	time.Sleep(s.Duration)
	log.Info().Msgf("done sleeping: %s", s.Duration)
	return nil
}
