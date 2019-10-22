package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/identifier"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	Context context.Context
	Cancel  context.CancelFunc
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())
}

func Configure() {
	flag.Parse()

	logging.Configure()
	identifier.Log()
}

func Wait() {
	<-Context.Done()
	log.Info().Msg("Context Done")
}
