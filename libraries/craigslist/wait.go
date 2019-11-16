package craigslist

import (
    "math/rand"
    "time"
    
    //"github.com/walterjwhite/go-application/libraries/logging"

	"github.com/rs/zerolog/log"
)

type RandomDelay struct {
	MinimumDelay int
	Deviation int
}

type FixedDelay struct {
	Delay int
}

type Waiter interface {
	Wait()
}

func (d *RandomDelay) Wait() {
	rand.Seed(time.Now().UnixNano())
    n := rand.Intn(d.Deviation) + d.MinimumDelay
	
	doWait(n)
}

func (d *FixedDelay) Wait() {
	doWait(d.Delay)
}
        
func doWait(durationInMillis int) {
	sleepTime := time.Duration(durationInMillis)*time.Millisecond
	
	log.Info().Msgf("sleeping %v", sleepTime)
	
	time.Sleep(sleepTime)
}
