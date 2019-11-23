// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"github.com/rs/zerolog/log"

	"context"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"github.com/walterjwhite/go-application/libraries/logging"
)

// approve post
//*[@id="new-edit"]/div/div[4]/div[1]/button
var ()

func init() {
	application.Configure()
}

func main() {
	actxt, _ /*cancelAllocator*/ := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222/devtools/browser/3fc4391e-e78a-4ede-bc0d-7cbfe8602507")
	//defer cancelAllocator()

	ctx /*cancelContext*/, _ := chromedp.NewContext(actxt) // create new tab
	//defer cancelContext()                            // close tab afterwards

	logging.Panic(chromedp.Run(ctx, chromedp.Navigate("https://fineuploader.com/demos")))

	var res interface{}
	logging.Panic(chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`window.alert = function (txt){return txt}`, &res)))

	//logging.Panic(chromedp.Run(ctx, chromedp.WaitReady("//*[@id=\"fine-uploader-gallery\"]/div/div[3]/input")))
	logging.Panic(chromedp.Run(ctx, chromedp.Click("//*[@id=\"fine-uploader-gallery\"]/div/div[3]/input")))

	log.Info().Msgf("Message: %v", res)

	/*
		printMsg := chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(lctx, func(ev interface{}) {

				if _, ok := ev.(*page.EventJavascriptDialogOpening); ok { // page loaded

					fmt.Printf(ev.(*page.EventJavascriptDialogOpening).Message) // holds msg!
				}
			}),
		}

		log.Info().Msgf("Message: %v", printMsg)
	*/
}
