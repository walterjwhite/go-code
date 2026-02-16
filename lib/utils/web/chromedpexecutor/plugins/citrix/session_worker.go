package citrix

import (
	"github.com/rs/zerolog/log"

	"sync"
)

func (s *Session) Work() {
	waitGroup := &sync.WaitGroup{}
	defer waitGroup.Wait()
	defer log.Info().Msg("session.Work() done")

	for index := range s.Instances {
		log.Debug().Msgf("[%v]", s.Instances[index])
		s.Instances[index].session = s

		waitGroup.Add(1)
		go s.Instances[index].run(waitGroup)
	}
}

func (s *Session) ShortBreak() {
	log.Info().Msgf("session.ShortBreak: %v", s.String())

	s.lockWorkers()
}

func (s *Session) LongBreak() {
	log.Info().Msgf("session.LongBreak: %v", s.String())

	s.lockWorkers()
}

func (s *Session) Lunch() {
	log.Info().Msgf("session.Lunch: %v", s.String())

	s.lockWorkers()
}

func (s *Session) Stop() {
	log.Info().Msgf("session.Stop: %v", s.String())

	s.cancel()
}
