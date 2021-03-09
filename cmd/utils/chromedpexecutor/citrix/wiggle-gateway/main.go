package main

import (
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/gateway"
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/gateway/cli"
	"time"
)

var (

	// TODO: randomize the interval, configure minimum interval and deviation ...
	tickleInterval = flag.String("TickleInterval", "3m", "Tickle Interval")
	session        = &gateway.Session{}
)

func init() {
	application.ConfigureWithProperties(session)

	i, err := time.ParseDuration(*tickleInterval)
	logging.Panic(err)

	session.Tickle = &gateway.Tickle{TickleInterval: &i}
}

func main() {
	session.Token = cli.New().Get()
	session.RunWith(application.Context, wiggleMouse)

	application.Wait()
}

func wiggleMouse() {
	for {
		session.ChromeDPSession.Execute(chromedp.KeyEvent(kb.Tab))
		time.Sleep(400 * time.Millisecond)
	}
}
