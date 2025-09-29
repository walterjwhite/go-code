package provider

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func New(c *Conf, pctx context.Context) (context.Context, context.CancelFunc) {
	ctx, _ := c.getAllocator(pctx)
	ctx, cancel := chromedp.NewContext(ctx /*, chromedp.WithDebugf(l.Printf)*/)

	ctx = c.withDelayer(ctx)

	if c.Headless {
		logging.Panic(chromedp.Run(ctx, chromedp.EmulateViewport(1920, 1080)))
	}

	logging.Panic(chromedp.Run(ctx))
	return ctx, cancel
}
