package gateway

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/periodic"
)

func (s *Session) tickle(ctx context.Context) {
	log.Info().Msg("tickle")

	if s.Tickle.periodicInstance != nil {
		log.Info().Msg("tickle exists")

		s.Tickle.periodicInstance.Cancel()
		s.Tickle.periodicInstance = nil
	}

	s.Tickle.periodicInstance = periodic.Periodic(ctx, s.Tickle.TickleInterval, s.doTickle)
	log.Info().Msgf("tickle instance: %v", s.Tickle.periodicInstance)
}

func (s *Session) doTickle() error {
	log.Info().Msgf("tickling: %v", s.Endpoint.Uri)
	s.chromedpsession.Execute(chromedp.Navigate(s.Endpoint.Uri))

	return nil
}
