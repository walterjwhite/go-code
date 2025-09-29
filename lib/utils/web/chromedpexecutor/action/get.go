package action

import (
	"context"

	"github.com/chromedp/chromedp"
)

func Get(ctx context.Context, selector interface{}) (string, error) {
	var value string
	err := chromedp.Run(ctx, chromedp.Text(selector, &value, chromedp.NodeVisible))

	return value, err
}
