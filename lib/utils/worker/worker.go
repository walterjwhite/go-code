package worker

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/until"
	"time"
)

func (c *Conf) Validate() error {
	if c.isPastEndTime() {
		return errors.New("end hour already passed")
	}

	return nil
}

func (c *Conf) Run(pctx context.Context) {
	log.Debug().Msgf("worker.Run.%s - start", c.worker.String())
	c.ctx, c.cancel = context.WithCancel(pctx)
	defer c.cancel()

	until.WaitUntil(c.StartHour, 0)

	for !c.isPastEndTime() {
		select {
		case <-c.ctx.Done():
			log.Warn().Msgf("worker.Run.%s - context is done", c.worker.String())
			return
		default:
			c.work()
			c.pomodoroCycles++

			c.takeBreak()
		}
	}

	log.Debug().Msgf("worker.Run.%s - past end time", c.worker.String())
	c.stop()
	log.Debug().Msgf("worker.Run.%s - end", c.worker.String())
}

func (c *Conf) isPastEndTime() bool {
	return time.Now().Hour() >= c.EndHour
}

func (c *Conf) takeBreak() {
	if c.pomodoroCycles%4 == 0 {
		c.longBreak()
	} else {
		c.shortBreak()
	}

	if c.isTimeForLunch() {
		c.lunch()
		c.hadLunch = true
	}
}

func (c *Conf) isTimeForLunch() bool {
	if c.hadLunch {
		log.Info().Msg("worker.isTimeForLunch - already had lunch")
		return false
	}

	now := time.Now()
	if now.Hour() < c.LunchStartHour {
		return false
	}

	log.Info().Msg("worker.isTimeForLunch - after lunch start")

	endLunchTime := time.Date(now.Year(), now.Month(), now.Day(),
		c.LunchStartHour, 0, 0, 0, now.Location())
	endLunchTime = endLunchTime.Add(c.LunchDuration)
	if now.Before(endLunchTime) {
		log.Info().Msg("worker.isTimeForLunch - before lunch end")
		return true
	}

	log.Info().Msg("worker.isTimeForLunch - lunchtime passed")
	c.hadLunch = true

	return false
}

func (c *Conf) Reset() {
	c.hadLunch = false
	c.pomodoroCycles = 0
}
