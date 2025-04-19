package action

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Location(ctx context.Context) string {
	var currentUrl string
	logging.Panic(chromedp.Run(ctx, chromedp.Location(&currentUrl)))

	return currentUrl
}
