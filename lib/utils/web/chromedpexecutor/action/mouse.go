package action

import (
	"context"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

type MouseLocation struct {
	X float64
	Y float64
}

const (
	mouseMoveEventListenerJavaScript = `
	if(globalThis.mouseListener !== null) {
		globalThis.mouseListener = true;

		document.addEventListener('mousemove', function(event) {
			window.mouseX = event.clientX;
			window.mouseY = event.clientY;
		});
	}
	`

	updateMouseLastPositionJavaScript = `
		window.lastMouseX = window.mouseX;
		window.lastMouseY = window.mouseY;
	`

	getMousePositionJavaScript     = `({x: window.mouseX, y: window.mouseY})`
	getLastMousePositionJavaScript = `({x: window.lastMouseX, y: window.lastMouseY})`

	delayBetweenMouseMoves = 1 * time.Second
)

var (
	wiggleDeltas [4][2]float64 = [4][2]float64{
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	}
)

func AttachMousePositionListener(ctx context.Context) error {
	log.Debug().Msg("attaching mouse position listener")

	return chromedp.Run(ctx,
		chromedp.Evaluate(mouseMoveEventListenerJavaScript, nil))
}

func GetMousePosition(ctx context.Context) (float64, float64, error) {
	var mouseLocation MouseLocation

	err := chromedp.Run(ctx,
		chromedp.Evaluate(getMousePositionJavaScript, &mouseLocation))
	if err != nil {
		return 0, 0, err
	}

	return mouseLocation.X, mouseLocation.Y, nil
}

func UpdateMousePosition(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.Evaluate(updateMouseLastPositionJavaScript, nil))
}

func WasMouseMoved(ctx context.Context) (error, bool) {
	mouseX, mouseY, err := GetMousePosition(ctx)
	if err != nil {
		logging.Warn(err, false, "GetMousePosition")
		return err, false
	}

	var lastMouseLocation MouseLocation
	err = chromedp.Run(ctx,
		chromedp.Evaluate(getLastMousePositionJavaScript, &lastMouseLocation))
	if err != nil {
		logging.Warn(err, false, "getLastMousePositionJavaScript")
		return err, false
	}

	moved := (mouseX != lastMouseLocation.X || mouseY != lastMouseLocation.Y)
	log.Debug().Msgf("mouse @: (%f, %f) <- (%f, %f)", mouseX, mouseY, lastMouseLocation.X, lastMouseLocation.Y)

	err = UpdateMousePosition(ctx)
	return err, moved
}

func MoveMouse(ctx context.Context, x, y float64) error {
	log.Debug().Msgf("moving mouse to: %f,%f", x, y)

	return chromedp.Run(ctx,
		chromedp.MouseEvent(input.MouseMoved, x, y),
	)
}

func Wiggle(ctx context.Context) error {
	err := AttachMousePositionListener(ctx)
	if err != nil {
		return err
	}

	x, y, err := GetMousePosition(ctx)
	if err != nil {
		return err
	}

	for i, delta := range wiggleDeltas {
		err, wasMoved := WasMouseMoved(ctx)
		if err != nil {
			return nil
		}
		if wasMoved {
			log.Warn().Msg("mouse was moved")
		}

		err = MoveMouse(ctx, x+delta[0], y+delta[1])
		if err != nil {
			return err
		}

		if i < len(wiggleDeltas)-1 {
			time.Sleep(delayBetweenMouseMoves)
		}
	}

	return nil
}
