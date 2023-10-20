package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func Click(ctx context.Context, visibleTimeout time.Duration, locateDelay delay.Delayer, selector interface{}, opts ...chromedp.QueryOption) error {
	err := Locate(ctx, visibleTimeout, locateDelay, selector, opts...)

	if err != nil {
		return err
	}

	return chromedp.Run(ctx, chromedp.Click(selector, opts...))
}
