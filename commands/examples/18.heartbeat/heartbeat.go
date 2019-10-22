package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/heartbeat"
	"time"
)

func main() {
	application.Configure()

	heartbeatAware()
}

func heartbeatAware() {
	heartbeat.Heartbeat(longRunningCall, heartbeater, 1*time.Second)
}

func longRunningCall() {
	time.Sleep(5 * time.Second)
	log.Info().Msg("Completed Execution")
}

func heartbeater() error {
	log.Debug().Msg("heartbeater")
	return nil
}
