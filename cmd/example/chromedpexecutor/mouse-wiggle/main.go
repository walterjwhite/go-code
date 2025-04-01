package main

import (
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"
	"time"
)

var (
	s session.ChromeDPSession
)

func init() {
	application.Configure()
	s = remote.New(application.Context)
}

func main() {
	defer s.Cancel()

	session.Execute(s, chromedp.Navigate("https://jsfiddle.net/m1erickson/WB7Zu/"), chromedp.Sleep(time.Duration(10*time.Second)))

	timeBetweenActions := time.Duration(1 * time.Second)

	for {
		log.Info().Msg("moving mouse to: 100,100")
		session.Execute(s,
			chromedp.MouseEvent(input.MouseMoved, 100, 100),
			chromedp.Sleep(timeBetweenActions))

		log.Info().Msg("moving mouse to: 200,100")
		session.Execute(s,
			chromedp.MouseEvent(input.MouseMoved, 200, 100),
			chromedp.Sleep(timeBetweenActions))

		log.Info().Msg("moving mouse to: 200,200")
		session.Execute(s,
			chromedp.MouseEvent(input.MouseMoved, 200, 200),
			chromedp.Sleep(timeBetweenActions))

		log.Info().Msg("moving mouse to: 100,200")
		session.Execute(s,
			chromedp.MouseEvent(input.MouseMoved, 100, 200),
			chromedp.Sleep(timeBetweenActions))
	}

}
