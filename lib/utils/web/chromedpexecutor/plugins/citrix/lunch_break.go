package citrix

import (
	"time"
)

func (s *Session) initLunchBreak(breakChannel chan *time.Duration) {
	if !waitUntil(s.LunchBreakStartHour) {
		return
	}

	t := time.Hour
	breakChannel <- &t
}
