package main

import (
	//"github.com/rs/zerolog/log"

	"context"

	"flag"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	devToolWsUrl = flag.String("devtools-ws-url", "", "DevTools Websocket URL")
)

func init() {
	application.Configure()
}

func main() {
	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), *devToolWsUrl)
	defer cancelActxt()

	ctx, cancelCtxt := chromedp.NewContext(actxt) // create new tab
	defer cancelCtxt()                            // close tab afterwards

	logging.Panic(chromedp.Run(ctx, chromedp.Navigate("https://ps.uci.edu/~franklin/doc/file_upload.html")))
	logging.Panic(chromedp.Run(ctx, chromedp.SendKeys("/html/body/form/input[1]", "/tmp/testing")))
	logging.Panic(chromedp.Run(ctx, chromedp.Click("/html/body/form/input[2]")))
}
