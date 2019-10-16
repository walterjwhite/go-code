package timeout

import (
	"context"
	"time"
)

type TimedExecution struct {
	function func()
}

type ContextAbortedException struct{}

func (c *ContextAbortedException) Error() string {
	return "Context was aborted"
}

func Limit(function func(), maximumExecutionTime time.Duration, parentContext context.Context) error {
	ctx, cancel := context.WithTimeout(parentContext, maximumExecutionTime)
	defer cancel()

	t := &TimedExecution{function: function}
	return t.call(ctx)
}

func (t *TimedExecution) call(ctx context.Context) error {
	invocationCompletedChannel := make(chan bool, 1)

	go func() {
		t.function()
		invocationCompletedChannel <- true
	}()

	select {

	case <-invocationCompletedChannel:
		return nil
	case <-ctx.Done():
		return &ContextAbortedException{}
	}
}
