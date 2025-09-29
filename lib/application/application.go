package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/property"

	"os"
	"os/signal"
	"syscall"
)

var (
	Context context.Context
	Cancel  context.CancelFunc
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())
	configureLogging()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan

		Cancel()
	}()
}

func Configure(configurations ...interface{}) {
	flag.Parse()
	Load(configurations...)

	doConfigure()
}

func Load(configurations ...interface{}) {
	for i := range configurations {
		property.Load(configurations[i])
	}
}

func doConfigure() {
	configureLogging()

	logIdentifier()
	logStart()
}

func logStart() {
	log.Info().Msg("Application started")
}

func Wait() {
	<-Context.Done()

	log.Info().Msg("Application Context Done")
}
