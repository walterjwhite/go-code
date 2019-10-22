package timeout

import (
	"context"
	"errors"
	"time"
)

type TimedExecution struct {
	function func()
}

func Limit(function func(), maximumExecutionTime time.Duration, parentContext context.Context) error {
	ctx, cancel := context.WithTimeout(parentContext, maximumExecutionTime)
	defer cancel()

	t := &TimedExecution{function: function}
	return t.call(ctx)
}

func (t *TimedExecution) call(ctx context.Context) error {
	invocationCompletedChannel := make(chan bool)

	go func() {
		t.function()
		invocationCompletedChannel <- true
	}()

	select {

	case <-invocationCompletedChannel:
		return nil
	case <-ctx.Done():
		return errors.New("Context was aborted")
	}
}
