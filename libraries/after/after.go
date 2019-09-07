package after

import (
	"context"
	"time"
)

type ErrorFunction struct {
	Function func() error
}

func After(ctx context.Context, delay time.Duration, fn func() error) *time.Timer {
	errorFunction := ErrorFunction{Function: fn}
	return time.AfterFunc(delay, errorFunction.function)
}

func (f *ErrorFunction) function() {
	err := f.Function()
	if err != nil {
		panic(err)
	}
}
