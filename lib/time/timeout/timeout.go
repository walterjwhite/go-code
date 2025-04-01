package timeout

import (
	"context"

	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

func Limit(function func(), maximumExecutionTime *time.Duration, parentContext context.Context) {
	ctx, cancel := context.WithTimeout(parentContext, *maximumExecutionTime)
	defer cancel()

	doLimit(ctx, function)
}

func doLimit(ctx context.Context, function func()) {
	invocationCompletedChannel := make(chan bool)

	go func() {
		function()
		invocationCompletedChannel <- true
	}()

	select {
	case <-invocationCompletedChannel:
		return
	case <-ctx.Done():
		logging.Panic(fmt.Errorf("context was aborted"))
	}
}
