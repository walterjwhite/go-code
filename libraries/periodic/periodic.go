package periodic

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func Periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// initial invocation
	logging.Panic(fn())

	for {
		select {
		case <-ticker.C:
			logging.Panic(fn())
		case <-ctx.Done():
			return
		}
	}
}

func GetInterval(intervalString string) time.Duration {
	duration, err := time.ParseDuration(intervalString)
	logging.Panic(err)

	return duration
}
