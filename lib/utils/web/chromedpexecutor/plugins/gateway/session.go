package gateway

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"

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
	time.Sleep(*s.PostAuthenticationDelay)

	if len(s.PostAuthenticationActions) > 0 {
		s.ChromeDPSession.Execute(run.ParseActions(s.PostAuthenticationActions...)...)
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
