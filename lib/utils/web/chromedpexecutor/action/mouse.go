package action

import (
	"context"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"time"
)

type MouseLocation struct {
	X             float64
	Y             float64
	ExpectedX     float64
	ExpectedY     float64

	MouseMoveTime int64
}

const (
	mouseMoveEventListenerJavaScript = `
	if (globalThis.mouseListener === undefined) {
		globalThis.mouseListener = true;
		document.addEventListener('mousemove', function(event) {
			if (!('expectedMouseX' in window)) {
				window.expectedMouseX = event.clientX;
				window.expectedMouseY = event.clientY;
			}

			window.mouseX = event.clientX;
			window.mouseY = event.clientY;
			window.mouseMoveTime = Math.floor(Date.now() / 1000);
		});
	}
	`

	getMousePositionJavaScript = `({
		x: window.mouseX,
		y: window.mouseY,

		expectedx: window.expectedMouseX,
		expectedy: window.expectedMouseY,
		
		mousemovetime: window.mouseMoveTime
	})`


	updateMousePositionJavaScript = `
		window.expectedMouseX = window.mouseX;
		window.expectedMouseY = window.mouseY;
	`


	delayBetweenMouseMoves = 1 * time.Second
	movementTimeout        = 30 * time.Second
)

var (
	wiggleDeltas [4][2]float64 = [4][2]float64{
		{10, 0},
		{10, 10},
		{0, 10},
		{0, 0},
	}
)

func updateMousePosition(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.Evaluate(updateMousePositionJavaScript, nil))
}

func AttachMousePositionListener(ctx context.Context) error {
	log.Debug().Msg("attaching mouse position listener")

	return chromedp.Run(ctx,
		chromedp.Evaluate(mouseMoveEventListenerJavaScript, nil))
}

func GetMousePosition(ctx context.Context) (MouseLocation, error) {
	var mouseLocation MouseLocation

	return mouseLocation, chromedp.Run(ctx,
		chromedp.Evaluate(getMousePositionJavaScript, &mouseLocation))
}

func WasMouseMoved(ctx context.Context) (bool, error) {
	mouseLocation, err := GetMousePosition(ctx)
	if err != nil {
		return false, err
	}

	lastMovementTime := time.Unix(mouseLocation.MouseMoveTime, 0)
	timeSinceLastMouseMovement := time.Since(lastMovementTime)

	log.Info().Msgf("WasMouseMoved.lastMovement: %f, %f [%f, %f] @ %v", mouseLocation.X, mouseLocation.Y, mouseLocation.ExpectedX, mouseLocation.ExpectedY, timeSinceLastMouseMovement)

	if mouseLocation.X == mouseLocation.ExpectedX &&
		mouseLocation.Y == mouseLocation.ExpectedY {
			return false, nil
		}
	
	wasMoved := timeSinceLastMouseMovement < movementTimeout
	err = updateMousePosition(ctx)

	return wasMoved, err
}

func MoveMouse(ctx context.Context, x, y float64) error {
	log.Debug().Msgf("moving mouse to: %f,%f", x, y)

	err := chromedp.Run(ctx, chromedp.MouseEvent(input.MouseMoved, x, y))
	if err != nil {
		return err
	}

	return updateMousePosition(ctx)
}

func Wiggle(ctx context.Context) error {
	err := AttachMousePositionListener(ctx)
	if err != nil {
		return err
	}

	mouseLocation, err := GetMousePosition(ctx)
	if err != nil {
		return err
	}

	for i, delta := range wiggleDeltas {
		wasMoved, err := WasMouseMoved(ctx)
		if err != nil {
			return nil
		}
		if wasMoved {
			log.Warn().Msg("mouse was moved")
		}

		err = MoveMouse(ctx, mouseLocation.X+delta[0], mouseLocation.Y+delta[1])
		if err != nil {
			return err
		}

		if i < len(wiggleDeltas)-1 {
			time.Sleep(delayBetweenMouseMoves)
		}
	}

	return nil
}
