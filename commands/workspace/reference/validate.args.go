package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func validateArgs() {
	if len(flag.Args()) < 2 {
		if flag.Args()[0] != "index" {
			logging.Panic(fmt.Errorf("Expecting action and reference names..."))
		}
	}
}
