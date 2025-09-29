package delay

import (
	"math/rand"
	"time"
)

type RandomDelay struct {
	amount time.Duration
}

func NewRandom(d time.Duration) *RandomDelay {
	return &RandomDelay{amount: d}
}

func (d *RandomDelay) Delay() {
	doDelay(time.Duration(rand.Int63n(d.amount.Nanoseconds())) * time.Nanosecond)
}
