package delay

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Delayer interface {
	Delay()
}

func doDelay(d time.Duration) {
	log.Debug().Msgf("sleeping %v", d)
	time.Sleep(d)
}
