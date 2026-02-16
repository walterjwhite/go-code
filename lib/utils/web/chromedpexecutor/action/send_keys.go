package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func SendKeys(ctx context.Context, visibleTimeout time.Duration, locateDelay delay.Delayer, selector interface{}, value string, opts ...chromedp.QueryOption) error {
	err := Locate(ctx, visibleTimeout, locateDelay, selector, opts...)

	if err != nil {
		return err
	}



	return chromedp.Run(ctx, chromedp.SendKeys(selector, value, opts...))
}
