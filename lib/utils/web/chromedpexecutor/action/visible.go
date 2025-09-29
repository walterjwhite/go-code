package action

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func IsVisible(ctx context.Context, selector interface{}, value string, opts ...chromedp.QueryOption) bool {
	vctx, vcancel := timeout(ctx, 100*time.Millisecond)
	defer vcancel()

	err := chromedp.Run(vctx, chromedp.WaitVisible(selector, opts...))
	return err == nil
}
