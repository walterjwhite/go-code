package headless

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func New(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx1, cancel := chromedp.NewContext(ctx)

	logging.Panic(chromedp.Run(ctx1, chromedp.EmulateViewport(1920, 1080)))

	return ctx1, cancel
}
