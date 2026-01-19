package main

import (
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"

	"time"
)

const (
	browserHeadersUrl = "https://www.whatismybrowser.com/detect/what-http-headers-is-my-browser-sending/"
	botDetectionUrl   = "https://fingerprint.com/products/bot-detection/"
)

func main() {
	ctx, cancel := provider.New(&provider.Conf{Headless: false}, application.Context)
	defer cancel()

	logging.Error(chromedp.Run(ctx,
		chromedp.Navigate(browserHeadersUrl),
		chromedp.WaitReady("body")))

	action.Screenshot(ctx, "/tmp/browser-headers.png")

	logging.Error(chromedp.Run(ctx,
		chromedp.Navigate(botDetectionUrl),
		chromedp.WaitReady("body")))

	action.Screenshot(ctx, "/tmp/bot-detection-url-fingerprint.com.png")

	time.Sleep(1 * time.Hour)
}
