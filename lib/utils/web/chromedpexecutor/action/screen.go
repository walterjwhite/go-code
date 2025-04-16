package action

import (
	"context"

	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

type Size struct {
	Width  int
	Height int
}

func GetScreenSize(ctx context.Context) Size {
	var screenSize Size
	err := chromedp.Run(ctx, chromedp.Evaluate(`({Width: screenWidth, Height: screenHeight})`, &screenSize))
	logging.Panic(err)

	return screenSize
}

func GetWindowSize(ctx context.Context) Size {
	var windowSize Size
	err := chromedp.Run(ctx, chromedp.Evaluate(`({Width: window.innerWidth, Height: window.innerHeight})`, &windowSize))
	logging.Panic(err)

	return windowSize
}
