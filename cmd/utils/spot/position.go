package main

import (
	"fmt"
	"github.com/walterjwhite/go/lib/external/spot"
)

func position(c *spot.Configuration) {
	fmt.Printf("latest position: %v, %v\n", c.Session.LatestReceivedRecord.Latitude, c.Session.LatestReceivedRecord.Longitude)
}
