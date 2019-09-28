package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/heartbeat"
	"log"
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
	log.Println("Completed Execution")
}

func heartbeater() error {
	log.Println("heartbeater")
	return nil
}
