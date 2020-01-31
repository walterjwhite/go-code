package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
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

	configureLogging()
	logIdentifier()
}

func Wait() {
	<-Context.Done()
	log.Info().Msg("Context Done")
}
