package movement

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/after"
	"time"
)

func (c *MovementConfiguration) schedule(duration time.Duration) {
	now := time.Now()

	start := getTime(now, c.StartHour, c.StartMinute)
	end := getTime(now, c.EndHour, c.EndMinute)

	if !withinWindow(now, start, end) {
		log.Debug().Msgf("Outside range: %v", now)
		return
	}

	// sends notification after movement window starts and duration has elapsed (if after is NOT canceled by a message update first)
	//u := time.Until(start.Add(c.MovementDurationTimeout))

	// log.Info().Msgf("Until: %v", u)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.after != nil {
		c.after.Cancel()
		c.after = nil
	}

	c.after = after.New(c.parentContext, &duration, c.onTimeout)
}

func getTime(now time.Time, hour, minute int) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(),
		hour,
		minute,
		0, 0,
		now.Location())
}

func withinWindow(now, start, end time.Time) bool {
	if now.Before(start) {
		return false
	}

	if now.After(end) {
		return false
	}

	return true
}
