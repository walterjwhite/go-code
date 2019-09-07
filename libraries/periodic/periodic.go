package periodic

import (
	"context"
	"time"
)

func Periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// initial invocation
	if err := fn(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
