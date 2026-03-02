package worker

import (
	"sync"
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

	mu sync.RWMutex
}

const (
	minHour = 0
	maxHour = 23
)

func (c *Conf) WithWorker(worker Worker) {
	if worker == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.worker = worker
}
