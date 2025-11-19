package action

import (
	"context"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

func OnTabClosed(ctx context.Context, fn func()) {
	chromedp.ListenBrowser(ctx, func(ev interface{}) {
		if ev, ok := ev.(*target.EventTargetDestroyed); ok {
			if c := chromedp.FromContext(ctx); c != nil {
				if c.Target.TargetID == ev.TargetID {
					log.Warn().Msg("OnTabClosed.detected tab/browser closure")

					select {
					case <-ctx.Done():
						log.Warn().Msg("OnTabClosed.ctx is done")
					default:
						log.Info().Msg("OnTabClosed.ctx is still alive")
					}

					fn()
				}
			}
		}
	})
}
