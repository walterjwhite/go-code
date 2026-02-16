package periodic

import (
	"context"
	"sync"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

type PeriodicInstance struct {
	function func() error

	ticker *time.Ticker

	ctx            context.Context
	cancelFunction context.CancelFunc

	mutex sync.Mutex
}

func Now(ctx context.Context, cancel context.CancelFunc, interval time.Duration, fn func() error) *PeriodicInstance {
	return Periodic(ctx, cancel, interval, true, fn)
}

func After(ctx context.Context, cancel context.CancelFunc, interval time.Duration, fn func() error) *PeriodicInstance {
	return Periodic(ctx, cancel, interval, false, fn)
}

func Periodic(ctx context.Context, cancel context.CancelFunc, interval time.Duration, runImmediately bool, fn func() error) *PeriodicInstance {
	ticker := time.NewTicker(interval)

	p := &PeriodicInstance{function: fn, ticker: ticker, ctx: ctx, cancelFunction: cancel}

	if runImmediately {
		p.doRun()
	}

	go p.run()
	go p.cancel()

	return p
}

func (p *PeriodicInstance) Cancel() {
	p.ticker.Stop()
	p.cancelFunction()
}

func (p *PeriodicInstance) doRun() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	logging.Warn(p.function(), "doRun")
}

func (p *PeriodicInstance) run() {
	for {
		select {
		case <-p.ticker.C:
			p.doRun()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *PeriodicInstance) cancel() {
	<-p.ctx.Done()
	p.ticker.Stop()
}

func (p *PeriodicInstance) Done() <-chan struct{} {
	return p.ctx.Done()
}
