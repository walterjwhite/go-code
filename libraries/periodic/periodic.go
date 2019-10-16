package periodic

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type PeriodicInstance struct {
	ticker *time.Ticker
}

func Periodic(ctx context.Context, interval time.Duration, fn func() error) *PeriodicInstance {
	ticker := time.NewTicker(interval)

	// initial invocation
	logging.Panic(fn())

	go run(fn, ticker)
	go cancel(ctx, ticker)

	return &PeriodicInstance{ticker: ticker}
}

func (p *PeriodicInstance) Cancel() {
	p.ticker.Stop()
}

func run(fn func() error, ticker *time.Ticker) {
	for {
		<-ticker.C
		logging.Panic(fn())
	}
}

func cancel(ctx context.Context, ticker *time.Ticker) {
	<-ctx.Done()
	ticker.Stop()
}

func GetInterval(intervalString string) time.Duration {
	duration, err := time.ParseDuration(intervalString)
	logging.Panic(err)

	return duration
}
