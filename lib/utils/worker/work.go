package worker

import (
	"github.com/rs/zerolog/log"
	"time"
)

func (c *Conf) work() {
	log.Debug().Msg("worker.work - start")
	ticker := time.NewTicker(c.WorkTickInterval)
	defer ticker.Stop()

	startTime := time.Now()

	if !c.doWork() {
		return
	}

	for range ticker.C {
		elapsed := time.Since(startTime)
		if elapsed >= c.WorkDuration {
			log.Debug().Msgf("worked for: %v - exiting", elapsed)
			return
		}

		if !c.doWork() {
			return
		}
	}

	log.Debug().Msg("worker.work() - end")
}

func (c *Conf) doWork() bool {
	select {
	case <-c.ctx.Done():
		log.Warn().Msg("worker.doWork - session interrupted!")
		return false
	default:
		log.Debug().Msg("worker.doWork - working")
		c.worker.Work()
	}

	return true
}

func (c *Conf) shortBreak() {
	log.Debug().Msg("worker.shortBreak - start")
	c.worker.ShortBreak()

	c.wait(c.ShortBreakDuration)

	log.Debug().Msg("worker.shortBreak - end")
}

func (c *Conf) longBreak() {
	log.Debug().Msg("worker.longBreak - start")
	c.worker.LongBreak()

	c.wait(c.LongBreakDuration)

	log.Debug().Msg("worker.longBreak - end")
}

func (c *Conf) lunch() {
	log.Debug().Msg("worker.lunch - start")
	c.worker.Lunch()

	c.wait(c.LunchDuration)

	log.Debug().Msg("worker.lunch - end")
}

func (c *Conf) wait(d time.Duration) {
	log.Debug().Msgf("waiting: %v", d)

	select {
	case <-time.After(d):
		log.Debug().Msgf("done waiting: %v", d)
	case <-c.ctx.Done():
		log.Info().Msg("worker.wait - context ended")
	}

	log.Debug().Msg("worker.wait - end")
}

func (c *Conf) stop() {
	log.Debug().Msg("worker.stop - start")

	c.worker.Stop()

	log.Debug().Msg("worker.stop - end")
}
