package timeout

import (
	"context"
	"fmt"
	"time"
)

type TimedExecution struct {
	MaximumExecutionTime time.Duration

	Function func()
}

func (t *TimedExecution) Error() string {
	return fmt.Sprintf("Invocation Timed Out after: %v\n", t.MaximumExecutionTime)
}

type ContextAbortedException struct{}

func (c *ContextAbortedException) Error() string {
	return "Context was aborted"
}

func Limit(function func(), maximumExecutionTime time.Duration, ctx context.Context) error {
	t := &TimedExecution{MaximumExecutionTime: maximumExecutionTime, Function: function}
	return t.call(ctx)
}

func (t *TimedExecution) call(ctx context.Context) error {
	c1 := make(chan bool, 1)

	go func() {
		t.Function()
		c1 <- true
	}()

	select {

	case <-c1:
		return nil
	case <-time.After(t.MaximumExecutionTime):
		return t
	case <-ctx.Done():
		return &ContextAbortedException{}
	}
}
