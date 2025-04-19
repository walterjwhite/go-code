package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"

	"github.com/rs/zerolog/log"
)

func Exists(pctx context.Context, d time.Duration, sel interface{}, opts ...chromedp.QueryOption) bool {
	log.Warn().Msgf("checking for visibility of: %s via %s", sel, pctx)

	ctx, cancel := context.WithTimeout(pctx, d)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Tasks{chromedp.WaitVisible(sel, opts...)})
	return err == nil
}
