package worker

import (
	"github.com/rs/zerolog/log"
)

func (c *Conf) stop() {
	if !c.EndTime.SleepUntil() {
		log.Warn().Msg("stop time already passed")
	}

	log.Info().Msg("time for quit")
	c.stopChannel <- true

}
