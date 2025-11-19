package mouse_driver

import (
	"time"
)

type Conf struct {
	MovementWaitTime   time.Duration
	TimeBetweenActions time.Duration
	Points             []Point // Change from Point[] to []Point
}

type Point struct {
	X float64
	Y float64
}

func New(movementWaitTime, timeBetweenActions time.Duration) *Conf {
	return &Conf{
		MovementWaitTime:   movementWaitTime,
		TimeBetweenActions: timeBetweenActions,
		Points: []Point{ // Use a slice instead of an array
			{X: 100, Y: 100},
			{X: 200, Y: 100},
			{X: 200, Y: 200},
			{X: 100, Y: 200},
		},
	}
}
