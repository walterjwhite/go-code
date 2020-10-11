package gateway

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/time/periodic"
)

func (s *Session) tickle(ctx context.Context) {
	log.Debug().Msg("tickle")

	if s.Tickle.periodicInstance != nil {
		log.Debug().Msg("tickle exists")

		s.Tickle.periodicInstance.Cancel()
		s.Tickle.periodicInstance = nil
	}

	// do not tickle immediately
	if s.Tickle != nil && s.Tickle.TickleInterval.Seconds() > 0 {
		s.Tickle.periodicInstance = periodic.Periodic(ctx, s.Tickle.TickleInterval, false, s.doTickle)
		log.Debug().Msgf("tickle instance: %v", s.Tickle.periodicInstance)
	} else {
		log.Debug().Msgf("not tickling: %v (seconds)", s.Tickle.TickleInterval.Seconds())
	}
}

func (s *Session) doTickle() error {
	log.Debug().Msgf("tickling: %v", s.Endpoint.Uri)
	s.ChromeDPSession.Execute(chromedp.Navigate(s.Endpoint.Uri))

	return nil
}
