package delay

import (
	"time"

	"github.com/rs/zerolog/log"
)

type DelayType int

const (
	Fixed DelayType = iota
	Random
)

type Delayer interface {
	Delay()
}

func doDelay(d time.Duration) {
	if d > 0 {
		log.Debug().Msgf("doDelay %v", d)
		time.Sleep(d)
		return
	}

	log.Debug().Msg("not sleeping, 0, specified")
}
