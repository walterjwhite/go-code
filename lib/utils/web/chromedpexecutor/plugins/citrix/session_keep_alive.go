package citrix

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	javascriptRefreshScript = `(
	function(){
		try {
			refreshButton = document.getElementsByClassName('messageBoxAction')[0].childNodes[0];
			if(refreshButton != null && refreshButton.offsetHeight > 0 && refreshButton.offsetWidth > 0) {
				refreshButton.click();
				return true;
			}
			
			return false;
		} catch(error) {
			return false;
		}
	}
	)()`
)

func (s *Session) keepAlive() {
	log.Debug().Msg("session.keepAlive - start")

	for range s.keepAliveTicker.C {
		log.Debug().Msg("session.keepAlive - checking if session is still active")
		if !s.isSessionStillActive() {
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

			logging.Warn(err, false, "keepAlive")
			if err != nil {
				log.Warn().Msg("session.keepAlive - error")
				s.cancel()
			}
		}
	}

	log.Debug().Msg("session.keepAlive - end")
}

func (s *Session) isSessionStillActive() bool {
	log.Debug().Msgf("session.isSessionStillActive: %v", s.Endpoint.Uri)

	select {
	case <-s.ctx.Done():
		return false
	default:
	}

	if isExpired(s.ctx) {
		log.Warn().Msg("session.isSessionStillActive - session appears to have expired (while running keep alive)")
		return false
	}

	log.Debug().Msg("session.isSessionStillActive - session appears to still be active")
	return true
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

func (s *Session) doTryKeepAliveByClickingRefreshButton(ctx context.Context) error {
	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(javascriptRefreshScript, &exists),
	)
	logging.Warn(err, false, "doTryKeepAlive.clickRefreshButton")

	log.Debug().Msgf("session.doTryKeepAlive.clickRefreshButton: %v", exists)

	logging.Warn(err, false, "doTryKeepAlive")
	return err
}
