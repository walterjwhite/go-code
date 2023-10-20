package delay

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Delayer interface {
	Delay()
}

func doDelay(d time.Duration) {
	if d > 0 {
		log.Debug().Msgf("sleeping %v", d)
		time.Sleep(d)
		return
	}

	log.Debug().Msg("not sleeping, 0, specified")
}
