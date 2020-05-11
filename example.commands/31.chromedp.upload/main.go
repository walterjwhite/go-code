package main

import (
	"github.com/rs/zerolog/log"

	"context"

	"flag"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"time"
)

var (
	devToolWsUrl = flag.String("devtools-ws-url", "", "DevTools Websocket URL")
)

func init() {
	application.Configure()
}

func main() {
	actxt, _ /*cancelActxt*/ := chromedp.NewRemoteAllocator(context.Background(), *devToolWsUrl)
	//defer cancelActxt()

	// create new tab
	ctx, _ /*cancelCtxt*/ := chromedp.NewContext(actxt)
	// close tab afterwards
	//defer cancelCtxt()

	//logging.Panic(chromedp.Run(ctx, upload()))
	log.Info().Msg("Navigating")
	logging.Panic(chromedp.Run(ctx, chromedp.Navigate("https://post.craigslist.org/k/0MDQTn8l6hGsbU3U-74MYw/6TZRk?s=editimage")))

	log.Info().Msg("Setting file")
	//logging.Panic(chromedp.Run(ctx, chromedp.SendKeys("//*[@id=\"uploader\"]/form/input[3]", []string{"/tmp/craigslist/San-Francisco.jpg", "/tmp/craigslist/San-Francisco-2.jpg"}, chromedp.NodeVisible)))
	logging.Panic(chromedp.Run(ctx, chromedp.SetUploadFiles("//*[@id=\"uploader\"]/form/input[3]", []string{"/tmp/craigslist/San-Francisco.jpg", "/tmp/craigslist/San-Francisco-2.jpg"}, chromedp.NodeVisible)))

	// 5 seconds / image
	time.Sleep(5 * time.Second)

	log.Info().Msg("Submitting")
	logging.Panic(chromedp.Run(ctx, chromedp.Click("/html/body/article/section/form/button")))
}
