// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	//"github.com/rs/zerolog/log"

	"context"

	"flag"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/craigslist"
	"github.com/walterjwhite/go-application/libraries/property"
	//"github.com/walterjwhite/go-application/libraries/logging"
)

// approve post
//*[@id="new-edit"]/div/div[4]/div[1]/button
var (
	devToolWsUrl   = flag.String("devtools-ws-url", "", "DevTools Websocket URL")
	craigslistPost = &craigslist.CraigslistPost{}
)

func init() {
	application.Configure()

	property.Load(craigslistPost, "")
}

func main() {
	e := chromedpexecutor.New(application.Context, *devToolWsUrl)
	defer e.Cancel()

	craigslistPost.Create(ctx)
}
