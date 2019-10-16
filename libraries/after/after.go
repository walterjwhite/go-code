package after

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type ErrorFunction struct {
	Function func() error
}

func After(ctx context.Context, delay time.Duration, fn func() error) *time.Timer {
	errorFunction := ErrorFunction{Function: fn}

	timer := time.AfterFunc(delay, errorFunction.function)
	go cancel(ctx, timer)

	return timer
}

func cancel(ctx context.Context, timer *time.Timer) {
	<-ctx.Done()
	timer.Stop()
}

func (f *ErrorFunction) function() {
	err := f.Function()
	logging.Panic(err)
}
