package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"log"
	"time"
)

func main() {
	application.Configure()

	timeConstrainedLongRunningCall()
}

func timeConstrainedLongRunningCall() {
	logging.Panic(timeout.Limit(longRunningCall, 3*time.Second))
}

func longRunningCall() {
	time.Sleep(5 * time.Second)
	log.Println("Completed Execution")
}
