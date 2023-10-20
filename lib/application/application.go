package application

import (
	"context"
	"flag"

	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/property"
	"github.com/walterjwhite/go-code/lib/application/shutdown"
)

var (
	Context context.Context
	Cancel  context.CancelFunc
	endCall sync.Once
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())

	configureLogging()
}

func Configure() {
	flag.Parse()
	doConfigure()
}

func ConfigureWithProperties(configurations ...interface{}) {
	for _, config := range configurations {
		property.Load(config)
	}

	doConfigure()
}

func doConfigure() {
	configureLogging()

	logIdentifier()
	logStart()
	shutdown.Add(Context, &defaultHandler{})
}

func logStart() {
	log.Info().Msg("Application started")
}

func OnEnd() {
	endCall.Do(doEnd)
}

type defaultHandler struct{}

func (a *defaultHandler) OnShutdown() {
	OnEnd()
}

func (a *defaultHandler) OnContextClosed() {
	OnEnd()
}

func doEnd() {
	log.Info().Msg("Application stopped")

	Cancel()
	os.Exit(0)
}

func Wait() {
	<-Context.Done()
	log.Info().Msg("Context Done")
}
