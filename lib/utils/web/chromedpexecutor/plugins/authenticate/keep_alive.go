package authenticate

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/time/wait"
)

func (s *Session) KeepAlive(ctx context.Context) {
	if /*s.IsKeepAlive*/ s.Website.SessionTimeout != nil && s.Website.KeepAliveUrl != nil {
		log.Debug().Msgf("running keep-alive: %v", s.Website.SessionTimeout)
		interval := *s.Website.SessionTimeout - 1*time.Minute
		wait.Wait(ctx, &interval, nil, s.onKeepAlive)
	}
}

func (s *Session) onKeepAlive() bool {
	s.chromedpsession.Execute(chromedp.Navigate(*s.Website.KeepAliveUrl))

	return false
}
