package action

import (
	"context"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type MouseLocation struct {
	X int
	Y int
}

func AttachMousePositionListener(ctx context.Context) {
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Evaluate(`document.addEventListener('mousemove', function(event) {
				window.mouseX = event.clientX;
				window.mouseY = event.clientY;
			});`, nil).Do(ctx)
		}))

	logging.Panic(err)
}

func GetMousePosition(ctx context.Context) (float64, float64) {
	var mouseLocation MouseLocation

	err := chromedp.Run(ctx,
		chromedp.Evaluate(`({x: window.mouseX, y: window.mouseY})`, &mouseLocation))
	logging.Panic(err)

	return float64(mouseLocation.X), float64(mouseLocation.Y)
}

func WasMouseMoved(ctx context.Context, lastX, lastY float64) (bool, float64, float64) {
	mouseX, mouseY := GetMousePosition(ctx)
	moved := (mouseX != lastX || mouseY != lastY)

	log.Info().Msgf("mouse @: %f, %f <- %f, %f", mouseX, mouseY, lastX, lastY)

	return moved, mouseX, mouseY
}

func MoveMouse(ctx context.Context, x, y float64) {
	log.Info().Msgf("moving mouse to: %f,%f", x, y)
	logging.Panic(chromedp.Run(ctx,
		chromedp.MouseEvent(input.MouseMoved, x, y)))
}

