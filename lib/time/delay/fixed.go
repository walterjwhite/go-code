package delay

import "time"

type FixedDelay struct {
	amount time.Duration
}

func New(d time.Duration) *FixedDelay {
	return &FixedDelay{amount: d}
}

func (d *FixedDelay) Delay() {
	doDelay(d.amount)
}
