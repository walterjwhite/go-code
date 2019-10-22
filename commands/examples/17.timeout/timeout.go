package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"time"
)

func main() {
	ctx := application.Configure()

	timeConstrainedLongRunningCall(ctx)
}

func timeConstrainedLongRunningCall(ctx context.Context) {
	logging.Panic(timeout.Limit(longRunningCall, 3*time.Second, ctx))
}

func longRunningCall() {
	time.Sleep(1 * time.Second)
	log.Info().Msg("Completed Execution")
}
