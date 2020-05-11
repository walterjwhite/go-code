package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/shutdown"
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
}

func logStart() {
	if !*noAuditFlag {
		log.Info().Msg("Application started")

		shutdown.Add(Context, &auditHandler{})
	}
}

// call via defer
func OnEnd() {
	go endCall.Do(doEnd)
}

type auditHandler struct{}

func (a *auditHandler) OnShutdown() {
	OnEnd()
}

func (a *auditHandler) OnContextClosed() {
	OnEnd()
}

func doEnd() {
	if !*noAuditFlag {
		log.Info().Msg("Application stopped")
	}

	Cancel()
}

func Wait() {
	<-Context.Done()
	log.Info().Msg("Context Done")
}
