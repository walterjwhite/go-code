package action

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Size struct {
	Width  int
	Height int
}

func GetScreenSize(ctx context.Context) (Size, error) {
	var screenSize Size
	err := chromedp.Run(ctx, chromedp.Evaluate(`({Width: window.screen.width, Height: window.screen.height})`, &screenSize))

	return screenSize, err
}

func GetWindowSize(ctx context.Context) (Size, error) {
	var windowSize Size
	err := chromedp.Run(ctx, chromedp.Evaluate(`({Width: window.innerWidth, Height: window.innerHeight})`, &windowSize))

	return windowSize, err
}
