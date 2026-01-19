package main

import (
	"log"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

var (
	browserConf = &provider.Conf{}
)

func init() {
	application.Configure(browserConf)
}

func main() {
	defer application.OnPanic()
	ctx, cancel := provider.New(browserConf, application.Context)
	defer cancel()

	var title string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://wikipedia.com/"),
		chromedp.Title(&title),
	); err != nil {
		log.Fatalf("Failed getting page's title: %v", err)
	}

	log.Println("Got title of:", title)
}
