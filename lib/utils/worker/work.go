package worker

import (
	"github.com/rs/zerolog/log"
	"time"
)

func (s *session) work() {
	log.Debug().Msg("worker.work - start")
	ticker := time.NewTicker(s.WorkTickInterval)
	defer ticker.Stop()

	startTime := time.Now()

	if !s.doWork() {
		return
	}

	for range ticker.C {
		elapsed := time.Since(startTime)
		if elapsed >= s.WorkDuration {
			log.Debug().Msgf("worked for: %v - exiting", elapsed)
			return
		}

		if !s.doWork() {
			return
		}
	}

	log.Debug().Msg("worker.work() - end")
}

func (s *session) doWork() bool {
	select {
	case <-s.ctx.Done():
		log.Warn().Msg("worker.doWork - session interrupted!")
		return false
	default:
		log.Debug().Msg("worker.doWork - working")
		s.worker.Work()
	}

	return true
}

func (s *session) shortBreak() {
	log.Debug().Msg("worker.shortBreak - start")
	s.worker.ShortBreak()

	s.wait(s.ShortBreakDuration)

	log.Debug().Msg("worker.shortBreak - end")
}

func (s *session) longBreak() {
	log.Debug().Msg("worker.longBreak - start")
	s.worker.LongBreak()

	s.wait(s.LongBreakDuration)

	log.Debug().Msg("worker.longBreak - end")
}

func (s *session) lunch() {
	log.Debug().Msg("worker.lunch - start")
	s.worker.Lunch()

	s.wait(s.LunchDuration)

	log.Debug().Msg("worker.lunch - end")
}

func (s *session) wait(d time.Duration) {
	log.Debug().Msgf("waiting: %v", d)

	select {
	case <-time.After(d):
		log.Debug().Msgf("done waiting: %v", d)
	case <-s.ctx.Done():
		log.Info().Msg("worker.wait - context ended")
	}

	log.Debug().Msg("worker.wait - end")
}

func (s *session) stop() {
	log.Debug().Msg("worker.stop - start")

	s.worker.Stop()

	log.Debug().Msg("worker.stop - end")
}
