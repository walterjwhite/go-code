package main

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/gen2brain/beeep"
)

func main() {
	// this appears to just do a terminal bell, not the hardware 'beep'
	logging.Panic(beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration))
}
