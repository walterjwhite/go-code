package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"time"
)

func init() {
	application.Configure()
}

func main() {
	timeConstrainedLongRunningCall()
}

func timeConstrainedLongRunningCall() {
	logging.Panic(timeout.Limit(longRunningCall, 3*time.Second, application.Context))
}

func longRunningCall() {
	time.Sleep(1 * time.Second)
	log.Info().Msg("Completed Execution")
}
