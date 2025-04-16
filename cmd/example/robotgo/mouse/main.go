package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

func main() {
	log.Info().Msg("Scroll -> 0,10")
	robotgo.Scroll(0, 10)

	log.Info().Msg("click -> left,true")
	robotgo.Click("left", true)

	log.Info().Msg("move 100,200 -> 1,1.0")
	robotgo.MoveSmooth(100, 200, 1.0, 1.0)


	robotgo.KeySleep = 100 // 100 millisecond

	for i := 0; i < 5; i++ {
		log.Info().Msgf("cmd+1 - %d", i)
		logging.Panic(robotgo.KeyTap("cmd", "1"))

		time.Sleep(1 * time.Second)
	}
}
