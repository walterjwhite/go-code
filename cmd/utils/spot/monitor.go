package main

import (
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/external/spot"
)

func doMonitor(c *spot.Configuration) {
	c.Monitor(application.Context)

	application.Wait()
}
