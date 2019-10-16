package main

import (
	"github.com/walterjwhite/go-application/libraries/after"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	ctx := application.Configure()

	a1 := after.After(ctx, 1*time.Second, afterOneSecond)
	a2 := after.After(ctx, 1*time.Minute, afterOneMinute)
	log.Info().Msg("Initialized timer")

	a1.Wait()

	log.Info().Msg("a1 Timer is complete, waiting on a2")
	a2.Cancel()
}

func afterOneSecond() error {
	log.Info().Msg("after 1 second has elapsed")
	return nil
}

func afterOneMinute() error {
	log.Info().Msg("after 1 minute has elapsed")
	return nil
}
