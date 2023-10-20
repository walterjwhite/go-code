package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func Locate(ctx context.Context, visibleTimeout time.Duration, locateDelay delay.Delayer, selector interface{}, opts ...chromedp.QueryOption) error {
	log.Debug().Msgf("timeout: %v", visibleTimeout)

	vctx, vcancel := timeout(ctx, visibleTimeout)
	defer vcancel()

	log.Debug().Msgf("set timeout: %#v", selector)

	err := chromedp.Run(vctx, chromedp.WaitVisible(selector, opts...))
	if err != nil {
		logging.Panic(err)
		return err
	}

	locateDelay.Delay()
	return chromedp.Run(ctx, chromedp.ScrollIntoView(selector, opts...))
}

func timeout(ctx context.Context, visibleTimeout time.Duration) (context.Context, context.CancelFunc) {
	if visibleTimeout > 0 {
		return context.WithTimeout(ctx, visibleTimeout)
	}

	return context.WithCancel(ctx)
}
