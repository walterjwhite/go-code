package delay

import (
	"math/rand"
	"time"
)

type RandomDelay struct {
	min       time.Duration
	deviation time.Duration
}

func NewRandom(m time.Duration, d time.Duration) *RandomDelay {
	return &RandomDelay{min: m, deviation: d}
}

func (d *RandomDelay) Delay() {
	if d.deviation > 0 {
		doDelay(d.min + time.Duration(rand.Int63n(d.deviation.Nanoseconds()))*time.Nanosecond)
		return
	}

	doDelay(d.min)
}
