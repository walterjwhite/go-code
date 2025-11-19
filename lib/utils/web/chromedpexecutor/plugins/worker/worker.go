package worker

import (
	"context"
)

type WorkerType int

const (
	MouseDriver WorkerType = iota
	Agent
	NOOP
)

func (w WorkerType) String() string {
	return [...]string{"MouseDriver", "Agent", "NOOP"}[w]
}

type ChromeDPWorker interface {
	Name() string
	Init(ctx context.Context, headless bool, contextuals ...interface{}) error
	Work(ctx context.Context, headless bool)
	Cleanup()
}
