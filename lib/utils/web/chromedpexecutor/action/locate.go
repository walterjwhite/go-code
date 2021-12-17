package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func Locate(ctx context.Context, visibleTimeout time.Duration, locateDelay delay.Delayer, selector interface{}, opts ...chromedp.QueryOption) error {
	vctx, vcancel := context.WithTimeout(ctx, visibleTimeout)
	defer vcancel()

	err := chromedp.Run(vctx, chromedp.WaitVisible(selector, opts...))
	if err != nil {
		return err
	}

	locateDelay.Delay()
	return chromedp.Run(ctx, chromedp.ScrollIntoView(selector, opts...))
}
