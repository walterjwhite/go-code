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

type session struct {
	*Conf
	ctx    context.Context
	cancel context.CancelFunc
	worker Worker
}

func (c *Conf) Run(pctx context.Context) {
	c.mu.RLock()
	worker := c.worker
	c.mu.RUnlock()

	workerString := worker.String()
	ctx, cancel := context.WithCancel(pctx)

	s := &session{
		Conf:   c,
		ctx:    ctx,
		cancel: cancel,
		worker: worker,
	}

	log.Debug().Msgf("worker.Run.%s - start", workerString)
	defer func() {
		s.cancel()
	}()

	until.WaitUntil(s.StartHour, 0)

	for !s.isPastEndTime() {
		select {
		case <-s.ctx.Done():
			log.Warn().Msgf("worker.Run.%s - context is done", workerString)
			return
		default:
			s.work()

			s.mu.Lock()
			s.pomodoroCycles++
			s.mu.Unlock()

			s.takeBreak()
		}
	}

	log.Debug().Msgf("worker.Run.%s - past end time", workerString)
	s.stop()
	log.Debug().Msgf("worker.Run.%s - end", workerString)
}

func (c *Conf) isPastEndTime() bool {
	return time.Now().Hour() >= c.EndHour
}

func (s *session) takeBreak() {
	s.mu.RLock()
	cycles := s.pomodoroCycles
	s.mu.RUnlock()

	if cycles%4 == 0 {
		s.longBreak()
	} else {
		s.shortBreak()
	}

	if s.isTimeForLunch() {
		s.lunch()
		s.mu.Lock()
		s.hadLunch = true
		s.mu.Unlock()
	}
}

func (s *session) isTimeForLunch() bool {
	s.mu.RLock()
	hadLunch := s.hadLunch
	s.mu.RUnlock()

	if hadLunch {
		log.Info().Msg("worker.isTimeForLunch - already had lunch")
		return false
	}

	now := time.Now()
	if now.Hour() < s.LunchStartHour {
		return false
	}

	log.Info().Msg("worker.isTimeForLunch - after lunch start")

	endLunchTime := time.Date(now.Year(), now.Month(), now.Day(),
		s.LunchStartHour, 0, 0, 0, now.Location())
	endLunchTime = endLunchTime.Add(s.LunchDuration)
	if now.Before(endLunchTime) {
		log.Info().Msg("worker.isTimeForLunch - before lunch end")
		return true
	}

	log.Info().Msg("worker.isTimeForLunch - lunchtime passed")
	s.mu.Lock()
	s.hadLunch = true
	s.mu.Unlock()

	return false
}

func (c *Conf) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hadLunch = false
	c.pomodoroCycles = 0
}
