package timeofday

import (
	"github.com/rs/zerolog/log"
	"time"
)

type TimeOfDay struct {
	Hour   int
	Minute int
}

func (t *TimeOfDay) Till() time.Duration {
	hours, minutes, seconds := time.Now().Clock()

	return time.Duration(t.Hour)*time.Hour + time.Duration(t.Minute)*time.Minute - (time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second)
}

func (t *TimeOfDay) SleepUntil() bool {
	sleepTime := t.Till()
	if sleepTime < 0 {
		log.Info().Msgf("not sleeping, time already passed: %v [%v]", t, sleepTime)
		return false
	}

	log.Info().Msgf("SleepUntil: %v [%v]", sleepTime, t)
	time.Sleep(sleepTime)

	return true
}
