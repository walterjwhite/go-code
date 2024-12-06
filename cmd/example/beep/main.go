package main

import (
	"github.com/gen2brain/beeep"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func main() {
	logging.Panic(beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration))
}
