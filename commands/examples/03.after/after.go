package main

import (
	"github.com/walterjwhite/go-application/libraries/after"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	ctx := application.Configure()

	after.After(ctx, 1*time.Second, afterOneSecond)
	//t2 := after.After(ctx, 1*time.Minute, afterOneMinute)
	log.Info().Msg("Initialized timer")

	//<-t1.C
	//<-t2.C
	log.Info().Msg("Timer is complete")
}

func afterOneSecond() error {
	log.Info().Msg("after 1 second has elapsed")
	return nil
}

func afterOneMinute() error {
	log.Info().Msg("after 1 minute has elapsed")
	return nil
}
