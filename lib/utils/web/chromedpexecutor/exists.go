package chromedpexecutor

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

func Exists(s session.ChromeDPSession, d time.Duration, sel interface{}, opts ...chromedp.QueryOption) bool {
	log.Warn().Msgf("checking for visibility of: %s via %s", sel, s)

	ctx, cancel := context.WithTimeout(s.Context(), d)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Tasks{chromedp.WaitVisible(sel, opts...)})
	return err == nil
}
