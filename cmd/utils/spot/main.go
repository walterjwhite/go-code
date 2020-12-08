package main

import (
	"flag"
	"fmt"

	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/external/spot"
)

var (
	feedId = flag.String("f", "", "Feed ID")
	action = flag.String("a", "monitor", "Action (monitor by default), export, position")
)

func init() {
	application.Configure()
}

func main() {
	validate()

	c := spot.New(*feedId)

	switch *action {
	case "monitor":
		doMonitor(c)
	case "position":
		position(c)
	case "export":
		export(c)
	default:
		logging.Panic(fmt.Errorf("-a, Action not understood %s", *action))
	}
}

func validate() {
	if len(*feedId) == 0 {
		logging.Panic(fmt.Errorf("-f, FeedID is required"))
	}

}
