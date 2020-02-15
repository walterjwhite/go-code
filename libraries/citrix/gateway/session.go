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

	// simply reload the page
	// TODO: block if tickling is active
	// with a 5-second delay after the fact before allowing subsequent execution
	s.tickle(ctx)

	// TODO: configure this
	time.Sleep(5 * time.Second)

	s.runPostAuthenticationActions(ctx)
}

func (s *Session) runPostAuthenticationActions(ctx context.Context) {
	/*
		s.launchDesktop()
		s.launchRemoteDesktop()
	*/

	if len(s.PostAuthenticationActions) > 0 {
		for i, a := range s.PostAuthenticationActions {
			if i > 0 {
				// TODO: configure this
				time.Sleep(5 * time.Second)

				log.Debug().Msgf("executing: %v", a.Name)

				for _, action := range a.Actions {
					s.chromedpsession.Execute(chromedpexecutor.GetScript(action))
				}
			}
		}
	}
}

func (s *Session) RunWith(ctx context.Context, fn func()) {
	s.Run(ctx)

	// after authenticated, run fn
	fn()
}

// TODO: configure this - post sign-in actions ...
func (s *Session) useLightVersion() {
	if s.UseLightVersion {
		s.chromedpsession.Execute(
			chromedp.Click("//*[@id=\"userMenuBtn\"]/p"),
			chromedp.Click("//*[@id=\"menuChangeClientBtn\"]"),
			chromedp.Click("//*[@id=\"changeclient-use-light-version\"]"),
		)
	}
}
