package after

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type AfterDelay struct {
	ctx   context.Context
	timer *time.Timer

	function func() error
	fired    chan bool
}

func After(ctx context.Context, delay time.Duration, fn func() error) *AfterDelay {
	afterDelay := &AfterDelay{ctx: ctx, function: fn, fired: make(chan bool)}
	afterDelay.timer = time.AfterFunc(delay, afterDelay.safeFunction)

	go afterDelay.onContextDone()

	return afterDelay
}

func (a *AfterDelay) Wait() {
	<-a.fired
}

func (a *AfterDelay) Cancel() {
	a.timer.Stop()
}

func (a *AfterDelay) onContextDone() {
	<-a.ctx.Done()
	a.Cancel()
}

func (a *AfterDelay) safeFunction() {
	logging.Panic(a.function())
}
