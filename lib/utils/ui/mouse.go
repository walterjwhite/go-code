package ui

import "github.com/go-vgo/robotgo"

var (
	lastX int
	lastY int
)

func WasMouseMoved() bool {
	currentX, currentY := robotgo.Location()
	defer updateLocation(currentX, currentY)

	return currentX != lastX || currentY != lastY
}

func updateLocation(x, y int) {
	lastX = x
	lastY = y
}
