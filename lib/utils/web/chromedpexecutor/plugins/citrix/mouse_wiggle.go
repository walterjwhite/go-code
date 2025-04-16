package citrix

import (
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

func (i *Instance) wiggleMouse() {
	log.Warn().Msgf("Mouse wiggle is enabled - %d", i.Index)

	if !i.session.Headless {
		action.AttachMousePositionListener(i.ctx)
	}

	i.handleExternalMouseMovement()

	i.moveMouse(100, 100)
	i.moveMouse(200, 100)
	i.moveMouse(200, 200)
	i.moveMouse(100, 200)

	i.captureScreenshot("/tmp/3.gateway-mouse-wiggle-%d.png")
}

func (i *Instance) handleExternalMouseMovement() {
	if !i.session.Headless {
		wasMoved, x, y := action.WasMouseMoved(i.ctx, i.lastMouseX, i.lastMouseY)
		if wasMoved {
			i.lastMouseX = x
			i.lastMouseY = y

			i.OnMouseMovedExternally()
			return
		}
	}
}

func (i *Instance) OnMouseMovedExternally() {
	log.Warn().Msgf("mouse was moved, waiting: %v", i.MovementWaitTime)
	time.Sleep(*i.MovementWaitTime)
}

func (i *Instance) moveMouse(x, y float64) {
	action.MoveMouse(i.ctx, x, y)
	i.lastMouseX = x
	i.lastMouseY = y
	time.Sleep(*i.TimeBetweenActions)
}
