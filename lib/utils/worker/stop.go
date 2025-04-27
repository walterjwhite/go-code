package worker

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
)

func (c *Conf) stop() {
	if !c.EndTime.SleepUntil() {
		log.Warn().Msg("stop time already passed")
	}

	log.Info().Msg("time to quit")
	c.stopChannel <- true

	application.Cancel()
}
