package main

import (
	"flag"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway/google"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway"

	"time"
)

var (

	tickleInterval = flag.String("i", "3m", "Tickle Interval")
	session        = &gateway.Session{}
	googleConf     = &google.Provider{}
)

func init() {
	application.ConfigureWithProperties(session)
	application.ConfigureWithProperties(googleConf)

	if len(*tickleInterval) > 0 {
		i, err := time.ParseDuration(*tickleInterval)
		logging.Panic(err)

		session.Tickle = &gateway.Tickle{TickleInterval: &i}
	}

	session.Validate()
	session.InitializeChromeDP(application.Context)
}

func main() {
	googleConf.ReadToken(session)
	session.KeepAlive(application.Context)

	application.Wait()
}
