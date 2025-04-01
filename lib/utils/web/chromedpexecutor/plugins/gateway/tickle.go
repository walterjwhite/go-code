package gateway

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"sync"
)

func (s *Session) KeepAlive(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	log.Debug().Msg("tickle")

	if s.Tickle.periodicInstance != nil {
		log.Debug().Msg("tickle exists")

		s.Tickle.periodicInstance.Cancel()
		s.Tickle.periodicInstance = nil
	}

	if s.Tickle != nil && s.Tickle.TickleInterval.Seconds() > 0 {
		s.Tickle.periodicInstance = periodic.Periodic(s.session.Context(), s.Tickle.TickleInterval, false, s.doKeepAlive)
		log.Debug().Msgf("tickle instance: %v", s.Tickle.periodicInstance)
	} else {
		log.Debug().Msgf("not tickling: %v (seconds)", s.Tickle.TickleInterval.Seconds())
	}
}

func (s *Session) doKeepAlive() error {
	log.Debug().Msgf("tickling: %v", s.Endpoint.Uri)
	session.Execute(s.session, chromedp.Navigate(s.Endpoint.Uri))

	return nil
}

