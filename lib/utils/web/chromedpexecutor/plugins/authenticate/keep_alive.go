package authenticate

import (
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/until"
)

func (s *Session) doKeepAlive() error {
	if len(s.Website.keepAliveActions) > 0 && s.Website.SessionTimeout != nil {
		log.Debug().Msgf("running keep-alive: %v", s.Website.SessionTimeout)
		interval := *s.Website.SessionTimeout - 1*time.Minute
		until.New(s.ctx, &interval, nil, s.onKeepAlive)
	}

	return nil
}

// TODO: only invoke this if we have been idle the length of the session timeout
// if we are actively doing stuff, we don't need to do anything
func (s *Session) onKeepAlive() bool {
	logging.Panic(chromedp.Run(s.chromedpsession.Context(), s.Website.keepAliveActions...))

	return false
}
