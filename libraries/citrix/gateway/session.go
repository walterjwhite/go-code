package gateway

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/logging"

	"github.com/rs/zerolog/log"
	"time"
)

// authenticate and keep the session alive ...
func (s *Session) Run(ctx context.Context) {
	s.Authenticate(ctx)

	time.Sleep(*s.Endpoint.AuthenticationDelay)

	if !s.isAuthenticated() {
		logging.Panic(fmt.Errorf("Unable to authenticate"))
	}

	s.useLightVersion()
	s.tickle(ctx)

	s.runPostAuthenticationActions(ctx)
}

func (s *Session) runPostAuthenticationActions(ctx context.Context) {
	// TODO: configure this
	//time.Sleep(5 * time.Second)

	if len(s.PostAuthenticationActions) > 0 {
		for i, a := range s.PostAuthenticationActions {
			if i > 0 {
				// TODO: configure this
				time.Sleep(5 * time.Second)

				log.Debug().Msgf("executing: %v", a.Name)

				s.ChromeDPSession.Execute(chromedpexecutor.ParseActions(a.Actions...)...)
			}
		}
	}
}

func (s *Session) RunWith(ctx context.Context, fn func()) {
	s.Run(ctx)

	// after authenticated, run fn
	fn()
}

const (
	menuChangeClientButtonXpath = "//*[@id=\"menuChangeClientBtn\"]"
	useLightVersionXpath        = "//*[@id=\"changeclient-use-light-version\"]"
)

// TODO: configure this - post sign-in actions ...
func (s *Session) useLightVersion() {
	if s.UseLightVersion {
		s.ChromeDPSession.Execute(
			chromedp.Click(menuButtonXpath),
			chromedp.Click(menuChangeClientButtonXpath),
			chromedp.Click(useLightVersionXpath),
		)
	}
}
