package wait

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/periodic"
	"github.com/walterjwhite/go-application/libraries/timeout"
)

type waitInstance struct {
	periodic       *periodic.PeriodicInstance
	timedExecution *timeout.TimedExecution
	function       func() bool

	channel chan bool
}

// calls the function periodically with the given interval until it returns true, the call times out, or the context is Done
func Wait(ctx context.Context, interval time.Duration, limit time.Duration, fn func() bool) {
	channel := make(chan bool)

	w := &waitInstance{channel: channel, function: fn}
	w.periodic = periodic.Periodic(ctx, interval, w.run)

	// wait until done
	logging.Panic(timeout.Limit(w.doWait, limit, ctx))
}

func (w *waitInstance) doWait() {
	<-w.channel
	close(w.channel)
}

func (w *waitInstance) cancel() {
	w.periodic.Cancel()
}

func (w *waitInstance) run() error {
	if w.function() {
		log.Info().Msg("Completed:")
		w.cancel()
		w.channel <- true

		return nil
	}

	log.Debug().Msg("Not yet completed:")
	return nil
}
