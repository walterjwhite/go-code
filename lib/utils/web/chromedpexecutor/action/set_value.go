package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func SetValue(ctx context.Context, visibleTimeout time.Duration, locateDelay delay.Delayer, selector interface{}, value string, opts ...chromedp.QueryOption) error {
	err := Locate(ctx, visibleTimeout, locateDelay, selector, opts...)

	if err != nil {
		return err
	}

	// log.Debug().Msgf("sending: %v", value)
	// if len(value) > 32 {

	// }

	return chromedp.Run(ctx, chromedp.Clear(selector), chromedp.SetValue(selector, value, opts...))
}
