package main

import (
	"github.com/walterjwhite/go-application/libraries/after"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/shutdown"

	"github.com/rs/zerolog/log"
	"time"
)

type afterFunction struct {
	message string
}

func init() {
	application.Configure()
}

func main() {
	a1f := &afterFunction{message: "after 1 minute has elapsed"}
	a := after.After(application.Context, 1*time.Second, a1f.afterPeriod)
	log.Debug().Msg("Initialized timer")

	shutdown.Add(&afterShutdownHandler{})

	//application.Wait()
	a.Wait()
}

func (a *afterFunction) afterPeriod() error {
	log.Info().Msg(a.message)
	return nil
}

type afterShutdownHandler struct{}

func (a *afterShutdownHandler) OnShutdown() {
	log.Info().Msg("On shutdown")
}

func (a *afterShutdownHandler) OnContextClosed() {
	log.Info().Msg("On context closed")
}
