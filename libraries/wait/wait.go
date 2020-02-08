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
	periodic *periodic.PeriodicInstance
	function func() bool

	channel chan bool
}

// calls the function periodically with the given interval until it returns true, the call times out, or the context is Done
func Wait(ctx context.Context, interval *time.Duration, limit *time.Duration, userFunction func() bool) {
	channel := make(chan bool)

	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w := &waitInstance{channel: channel, function: userFunction}
	w.periodic = periodic.Periodic(wctx, interval, w.monitorFunction)

	// wait until done
	logging.Panic(timeout.Limit(w.doWait, limit, wctx))
}

func (w *waitInstance) doWait() {
	<-w.channel
	close(w.channel)
}

func (w *waitInstance) cancel() {
	w.periodic.Cancel()
}

func (w *waitInstance) monitorFunction() error {
	if w.function() {
		log.Info().Msg("Completed:")
		w.channel <- true
		w.cancel()

		return nil
	}

	log.Debug().Msg("Not yet completed:")
	return nil
}
