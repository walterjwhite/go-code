package sleep

import (
	"time"

	"github.com/walterjwhite/go-code/lib/security/random"

	"github.com/rs/zerolog/log"
)

type RandomDelay struct {
	MinimumDelay int
	Deviation    int
}

type FixedDelay struct {
	Delay int
}

type Waiter interface {
	Wait()
}

func (d *RandomDelay) Wait() {
	doWait(random.Of(d.Deviation) + d.MinimumDelay)
}

func (d *FixedDelay) Wait() {
	doWait(d.Delay)
}

func doWait(durationInMillis int) {
	sleepTime := time.Duration(durationInMillis) * time.Millisecond

	log.Info().Msgf("sleeping %v", sleepTime)

	time.Sleep(sleepTime)
}
