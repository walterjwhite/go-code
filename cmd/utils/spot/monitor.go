package main

import (
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/external/spot"
)

// TODO: configure actions (export of daily GPS data, export of last GPS track, email, sms, etc)
// standard (record GPS data to new file each day)
func doMonitor(c *spot.Configuration) {
	c.Monitor(application.Context)

	application.Wait()
}
