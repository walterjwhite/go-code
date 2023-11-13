package main

import (
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway/cli"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"time"
)

var (

	tickleInterval = flag.String("TickleInterval", "3m", "Tickle Interval")
	s              = &gateway.Session{}
)

func init() {
	application.ConfigureWithProperties(s)

	i, err := time.ParseDuration(*tickleInterval)
	logging.Panic(err)

	s.Tickle = &gateway.Tickle{TickleInterval: &i}
}

func main() {
	s.Token = cli.New().Get()
	s.RunWith(application.Context, wiggleMouse)

	application.Wait()
}

func wiggleMouse() {
	for {
		session.Execute(s.Session(), chromedp.KeyEvent(kb.Tab))
		time.Sleep(400 * time.Millisecond)
	}
}
