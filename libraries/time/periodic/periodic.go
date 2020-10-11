package periodic

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"sync"
	"time"
)

type PeriodicInstance struct {
	function func() error

	ticker *time.Ticker

	ctx            context.Context
	cancelFunction context.CancelFunc
	mutex          *sync.RWMutex
	runCount       int
}

func Now(parentContext context.Context, interval *time.Duration, fn func() error) *PeriodicInstance {
	return Periodic(parentContext, interval, true, fn)
}

func Periodic(parentContext context.Context, interval *time.Duration, runImmediately bool, fn func() error) *PeriodicInstance {
	ticker := time.NewTicker(*interval)

	ctx, cancel := context.WithCancel(parentContext)

	p := &PeriodicInstance{function: fn, ticker: ticker, ctx: ctx, cancelFunction: cancel, mutex: &sync.RWMutex{}}

	// initial invocation
	p.tryRun()

	go p.run()
	go p.cancel()

	return p
}

func (p *PeriodicInstance) Cancel() {
	p.ticker.Stop()
}

func (p *PeriodicInstance) run() {
	for {
		<-p.ticker.C
		p.tryRun()
	}
}

func (p *PeriodicInstance) tryRun() {
	p.mutex.RLock()
	count := p.runCount
	p.mutex.RUnlock()

	if count == 0 {
		p.doRun()
	}
}

func (p *PeriodicInstance) doRun() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.runCount++

	logging.Panic(p.function())
	p.runCount--
}

func (p *PeriodicInstance) cancel() {
	<-p.ctx.Done()
	p.ticker.Stop()

	p.cancelFunction()
}

func GetInterval(intervalString string) *time.Duration {
	duration, err := time.ParseDuration(intervalString)
	logging.Panic(err)

	return &duration
}
