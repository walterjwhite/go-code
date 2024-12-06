package until

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/periodic"
)

type instance struct {
	periodic *periodic.PeriodicInstance
	function func() bool

	channel chan bool
}

func New(ctx context.Context, interval *time.Duration, limit *time.Duration, f func() bool) {
	channel := make(chan bool, 1)

	ctx, cancel := context.WithCancel(ctx)
	if limit != nil {
		ctx, cancel = context.WithTimeout(ctx, *limit)
	}

	defer cancel()

	w := &instance{channel: channel, function: f}
	w.periodic = periodic.Now(ctx, interval, w.monitorFunction)

	w.doWait()
}

func (w *instance) doWait() {
	defer close(w.channel)

	<-w.channel
}

func (w *instance) Cancel() {
	w.periodic.Cancel()
}

func (w *instance) monitorFunction() error {
	if w.function() {
		log.Debug().Msg("Completed:")
		w.channel <- true

		return nil
	}

	log.Debug().Msg("Not yet completed:")
	return nil
}
