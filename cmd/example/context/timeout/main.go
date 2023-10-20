package main

import (
	"time"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/time/timeout"
)

func init() {
	application.Configure()
}

func main() {
	// _, cancel := context.WithTimeout(application.Context, 1*time.Second)
	// defer cancel()

	// log.Info().Msg("Attempting to sleep for 5 seconds, with a timeout of 1 second")
	// time.Sleep(5 * time.Second)
	// log.Error().Msg("context should have aborted sleeping after 1 second, how did we get here?")

	d := 1 * time.Second
	timeout.Limit(sleep, &d, application.Context)
}

func sleep() {
	time.Sleep(5 * time.Second)
}
