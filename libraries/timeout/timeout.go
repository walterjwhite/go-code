package timeout

import (
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

func Limit(function func(), maximumExecutionTime time.Duration) error {
	t := &TimedExecution{MaximumExecutionTime: maximumExecutionTime, Function: function}
	return t.Call()
}

func (t *TimedExecution) Call() error {
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
	}
}
