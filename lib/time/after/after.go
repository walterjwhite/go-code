package after

import (
	"context"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

type AfterDelay struct {
	ctx    context.Context
	cancel context.CancelFunc
	timer  *time.Timer

	function func() error
}

func New(ctx context.Context, delay *time.Duration, fn func() error) *AfterDelay {
	acontext, acancel := context.WithCancel(ctx)

	afterDelay := &AfterDelay{ctx: acontext, cancel: acancel, function: fn}
	afterDelay.timer = time.AfterFunc(*delay, afterDelay.safeFunction)

	return afterDelay
}

func (a *AfterDelay) Wait() {
	<-a.ctx.Done()
}

func (a *AfterDelay) Cancel() {
	defer a.cancel()

	a.timer.Stop()
}

func (a *AfterDelay) safeFunction() {
	defer a.cancel()

	logging.Panic(a.function())
}
