package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	
	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
	
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/gateway"
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/gateway/cli"
	"time"
)

var (

	// TODO: randomize the interval, configure minimum interval and deviation ...
	tickleInterval = flag.String("i", "", "Tickle Interval, disabled")
	session        = &gateway.Session{}
)

func init() {
	application.ConfigureWithProperties(session)

	if len(*tickleInterval) > 0 {
		i, err := time.ParseDuration(*tickleInterval)
		logging.Panic(err)

		session.Tickle = &gateway.Tickle{TickleInterval: &i}
	}
}

func main() {
	session.Token = cli.New().Get()

	log.Debug().Msgf("token: %v", session.Token)

	session.Run(application.Context)

	application.Wait()
}
