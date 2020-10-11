package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application/shutdown"
	"os"
	"sync"
)

var (
	Context context.Context
	Cancel  context.CancelFunc
	endCall sync.Once

	noAuditFlag = flag.Bool("NoAudit", false, "Disable Audit execution")
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())
}

func Configure() {
	flag.Parse()

	configureLogging()
	logIdentifier()

	logStart()
	shutdown.Add(Context, &defaultHandler{})
}

func logStart() {
	if !*noAuditFlag {
		log.Info().Msg("Application started")
	}
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
	if !*noAuditFlag {
		log.Info().Msg("Application stopped")
	}

	Cancel()
	os.Exit(0)
}

func Wait() {
	<-Context.Done()
	log.Info().Msg("Context Done")
}
