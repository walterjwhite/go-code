package citrix

import (
	"context"

	"github.com/chromedp/cdproto/target"

	"github.com/rs/zerolog/log"

	"strings"
	"sync"
)

const (
	tabType = "page"
)

func (s *Session) Work() {
	s.waitGroup = &sync.WaitGroup{}

	for index := range s.Instances {
		log.Debug().Msgf("[%v]", s.Instances[index])
		s.Instances[index].session = s

		s.waitGroup.Add(1)
		go s.Instances[index].run()
	}

	s.waitGroup.Wait()
	log.Info().Msg("session.Work() done")
}

func (s *Session) ShortBreak() {
	log.Info().Msgf("session.ShortBreak: %v", s.Name())
	s.lockWorkers()
}

func (s *Session) LongBreak() {
	log.Info().Msgf("session.LongBreak: %v", s.Name())
	s.lockWorkers()
}

func (s *Session) Lunch() {
	log.Info().Msgf("session.Lunch: %v", s.Name())
	s.lockWorkers()
}

func (s *Session) Stop() {
	log.Info().Msgf("session.Stop: %v", s.Name())
	s.lockWorkers()
}

type CitrixWorker interface {
	Name() string
	Work(ctx context.Context, headless bool)
	Cleanup()
}

func (s *Session) Name() string {
	var builder strings.Builder

	builder.WriteString("citrix: {")

	for index := range s.Instances {
		if s.Instances[index].Worker != nil {
			builder.WriteString(" worker: " + s.Instances[index].Worker.Name())
		} else {
			builder.WriteString(" worker: NOT INITIALIZED")
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}

func matchTabWithNonEmptyURL(info *target.Info) bool {
	return info.Type == tabType && info.URL != ""
}
