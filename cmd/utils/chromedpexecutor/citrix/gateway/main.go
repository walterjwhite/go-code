package main

import (
	"flag"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/gateway/cli"
	"time"
)

var (

	tickleInterval = flag.String("i", "3m", "Tickle Interval")
	session        = &gateway.Session{}
)

func init() {
	application.ConfigureWithProperties(session)

	if len(*tickleInterval) > 0 {
		i, err := time.ParseDuration(*tickleInterval)
		logging.Panic(err)

		session.Tickle = &gateway.Tickle{TickleInterval: &i}
	}

	session.Validate()
	session.InitializeChromeDP(application.Context)
}

func main() {
	token := cli.New().Get()
	if !session.Run(token) {
		logging.Panic(fmt.Errorf("Unable to authenticate"))
	}

	session.KeepAlive(application.Context)
	application.Wait()
}
