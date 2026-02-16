package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"time"
)

func ScrollToEnd(iterationDelay time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		for {
			err := chromedp.Run(ctx, chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil))
			if err != nil {
				return err
			}

			time.Sleep(iterationDelay)

			var scrollHeight int64
			var clientHeight int64
			err = chromedp.Run(ctx, chromedp.Evaluate(`document.body.scrollHeight`, &scrollHeight))
			if err != nil {
				return err
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.body.clientHeight`, &clientHeight))
			if err != nil {
				return err
			}

			if scrollHeight <= clientHeight {
				break
			}
		}
		return nil
	})
}

func End(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.KeyEvent(kb.End))
}

func EndAction() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.KeyEvent(kb.End))
	})
}
