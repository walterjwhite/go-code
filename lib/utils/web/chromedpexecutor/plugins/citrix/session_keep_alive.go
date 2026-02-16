package citrix

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

func (s *Session) keepAlive() {
	keepAliveTicker := time.NewTicker(*s.Timeout)
	defer keepAliveTicker.Stop()

	log.Info().Msg("session.keepAlive - start")

	for range keepAliveTicker.C {
		log.Debug().Msg("session.keepAlive - checking if session is still active")
		if s.IsExpired() {
			log.Warn().Msg("session.keepAlive - IsExpired -> true")
			return
		}

		err := retry.Do(
			func() error {
				return s.doTryKeepAlive()
			},
			retry.Attempts(s.KeepAliveTries),
			retry.Delay(*s.KeepAliveDelay),
		)

		if err == nil {
			action.Screenshot(s.ctx, "/tmp/citrix-keep-alive.png")
		} else {
			action.Screenshot(s.ctx, "/tmp/citrix-keep-alive-timeout-error.png")
			if err != nil {
				log.Warn().Msg("session.keepAlive - error")
				return
			}
		}
	}

	log.Info().Msg("session.keepAlive - end")
}

func (s *Session) IsExpired() bool {
	log.Debug().Msgf("session.IsExpired: %v", s.Endpoint.Uri)

	select {
	case <-s.ctx.Done():
		return true
	default:
	}

	if IsContextExpired(s.ctx) {
		log.Warn().Msg("session.IsExpired - session appears to have expired (while running keep alive)")
		return true
	}

	log.Debug().Msg("session.IsExpired - session appears to still be active")
	return false
}

func (s *Session) doTryKeepAlive() error {
	log.Debug().Msg("session.doTryKeepAlive - start")

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}

	ctx, cancel := context.WithTimeout(s.ctx, *s.KeepAliveTimeout)
	defer cancel()

	return s.doTryKeepAliveByRefresh(ctx)
}

func (s *Session) doTryKeepAliveByRefresh(ctx context.Context) error {
	log.Debug().Msgf("session.doTryKeepAlive: %v", s.Endpoint.Uri)
	return chromedp.Run(ctx, chromedp.Navigate(s.Endpoint.Uri))
}
