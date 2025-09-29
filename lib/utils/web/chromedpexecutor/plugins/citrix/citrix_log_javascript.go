package citrix

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *Session) captureJavascript() context.CancelFunc {
	err := chromedp.Run(s.ctx, network.Enable())
	if err != nil {
		logging.Warn(err, false, "Session.captureJavascript - error enabling network")
		return nil
	}

	listenCtx, listenCancel := context.WithCancel(s.ctx)
	chromedp.ListenTarget(listenCtx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			if e.Type == network.ResourceTypeScript && e.Request.URL != "" {
				log.Debug().Msgf("Session.captureJavascript - javascript request: %s", e.Request.URL)
			}
		case *network.EventResponseReceived:
			if e.Type == network.ResourceTypeScript && e.Response.URL != "" {
				log.Debug().Msgf("Session.captureJavascript - javascript response: %s", e.Response.URL)
			}
		}
	})

	return listenCancel
}
