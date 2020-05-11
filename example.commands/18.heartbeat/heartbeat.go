package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/heartbeat"
	"time"
)

func init() {
	application.Configure()
}

func main() {
	heartbeatAware()
}

func heartbeatAware() {
	t := 1 * time.Second
	heartbeat.Heartbeat(longRunningCall, heartbeater, &t)
}

func longRunningCall() {
	time.Sleep(5 * time.Second)
	log.Info().Msg("Completed Execution")
}

func heartbeater() error {
	log.Debug().Msg("heartbeater")
	return nil
}
