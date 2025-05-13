package mouse_wiggle

import (
  "context"
  "github.com/rs/zerolog/log"

  "github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

  "time"
)

func (i *State) Name() string {
  return "mouse wiggle"
}

func (i *State) Work(ctx context.Context, headless bool) {
  log.Debug().Msgf("Mouse wiggle is enabled: %s", action.Location(ctx))

  if !headless {
    if !i.initialized {
      action.AttachMousePositionListener(ctx)
      i.initialized = true
    }

    if i.handleExternalMouseMovement(ctx) {
      return
    }
  }

  i.moveMouse(ctx, 100, 100)
  i.moveMouse(ctx, 200, 100)
  i.moveMouse(ctx, 200, 200)
  i.moveMouse(ctx, 100, 200)

  action.FullScreenshot(ctx, "/tmp/3.gateway-mouse-wiggle.png")
}

func (i *State) handleExternalMouseMovement(ctx context.Context) bool {
  wasMoved, x, y := action.WasMouseMoved(ctx, i.lastMouseX, i.lastMouseY)
  if wasMoved {
    i.lastMouseX = x
    i.lastMouseY = y

    i.onMouseMovedExternally()
    return true
  }

  return false
}

func (i *State) onMouseMovedExternally() {
  log.Warn().Msg("mouse was moved, skipping run")
}

func (i *State) moveMouse(ctx context.Context, x, y float64) {
  action.MoveMouse(ctx, x, y)
  i.lastMouseX = x
  i.lastMouseY = y
  time.Sleep(*i.TimeBetweenActions)
}

func (i *State) Cleanup() {
  i.initialized = false
}
