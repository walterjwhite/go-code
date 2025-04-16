package citrix

import (
	"github.com/rs/zerolog/log"
	"time"
)

type PomodoroInstance struct {
	OnDuration    *time.Duration
	OffDuration   *time.Duration
	BreakDuration *time.Duration

	cycle int
}

func (p *PomodoroInstance) init(breakChannel chan *time.Duration) {
	if p == nil || p.OnDuration == nil {
		log.Warn().Msg("disabling taking breaks")
		return
	}

	for {
		time.Sleep(*p.OnDuration)

		if p.cycle%4 == 0 {
			breakChannel <- p.BreakDuration
			time.Sleep(*p.BreakDuration)
		} else {
			breakChannel <- p.OffDuration
			time.Sleep(*p.OffDuration)
		}

		p.cycle++
	}
}
