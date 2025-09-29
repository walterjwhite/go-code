package worker

import (
	"context"
	"time"
)

type Conf struct {
	WorkDuration     time.Duration
	WorkTickInterval time.Duration

	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration

	LunchStartHour int
	LunchDuration  time.Duration

	StartHour int
	EndHour   int

	hadLunch       bool
	pomodoroCycles int

	worker Worker
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *Conf) WithWorker(worker Worker) {
	c.worker = worker
}
