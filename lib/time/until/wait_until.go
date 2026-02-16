package until

import (
	"github.com/rs/zerolog/log"
	"time"
)

func WaitUntil(targetHour, targetMinute int) {
	now := time.Now()

	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, now.Location())

	if targetTime.Before(now) {
		log.Warn().Msgf("time: %v has already passed", targetTime)
		return
	}

	waitDuration := targetTime.Sub(now)

	log.Info().Msgf("Waiting until %v", targetTime)
	time.Sleep(waitDuration)

	log.Info().Msg("done waiting")
}
