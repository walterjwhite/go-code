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


	d := 1 * time.Second
	timeout.Limit(sleep, &d, application.Context)
}

func sleep() {
	time.Sleep(5 * time.Second)
}
