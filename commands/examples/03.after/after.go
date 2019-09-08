package main

import (
	"github.com/walterjwhite/go-application/libraries/after"
	"github.com/walterjwhite/go-application/libraries/application"

	"log"
	"time"
)

func main() {
	ctx := application.Configure()

	timer := after.After(ctx, 1*time.Second, onAfter)
	log.Println("Initialized timer")

	<-timer.C
	log.Println("Timer is complete")
}

func onAfter() error {
	log.Println("after 1 second has elapsed")
	return nil
}
