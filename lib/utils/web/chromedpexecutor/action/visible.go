package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func IsVisible(pctx context.Context, selector interface{}, value string, opts ...chromedp.QueryOption) bool {
	ctx, cancel := timeout(pctx, 100*time.Millisecond)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.WaitVisible(selector, opts...))
	return err == nil
}
