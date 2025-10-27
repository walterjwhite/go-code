package worker

import (
	"context"
)

type WorkerType int

const (
	MouseWiggler WorkerType = iota
	Agent
	NOOP
)

func (w WorkerType) String() string {
	return [...]string{"MouseWiggler", "Agent", "NOOP"}[w]
}

type ChromeDPWorker interface {
	Name() string
	Init(ctx context.Context, headless bool, contextuals ...interface{}) error
	Work(ctx context.Context, headless bool)
	Cleanup()
}
