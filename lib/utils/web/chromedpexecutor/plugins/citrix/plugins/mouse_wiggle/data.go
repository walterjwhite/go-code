package mouse_wiggle

import (
	"time"
)

type State struct {
	MovementWaitTime   *time.Duration
	TimeBetweenActions *time.Duration

	lastMouseX float64
	lastMouseY float64

	initialized bool
}
