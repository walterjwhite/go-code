package gateway

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"

	"time"
)

const (
	menuChangeClientButtonXpath = "//*[@id=\"menuChangeClientBtn\"]"
	useLightVersionXpath        = "//*[@id=\"changeclient-use-light-version\"]"

	useLightVersionPromptXpath = "//*[@id=\"protocolhandler-welcome-useLightVersionLink\"]"
)

// authenticate and keep the session alive ...
func (s *Session) Run(token string) bool {
	log.Info().Msgf("running with: %v", token)
	validateToken(token)

	s.Authenticate(token)

	if !s.IsAuthenticated() {
		return false
	}

	s.useLightVersion()
	s.runPostAuthenticationActions()

	return true
}

func (s *Session) runPostAuthenticationActions() {
	if len(s.PostAuthenticationActions) > 0 {
		log.Info().Msgf("running post authentication actions - delay: %v", *s.PostAuthenticationDelay)
		time.Sleep(*s.PostAuthenticationDelay)

		log.Info().Msgf("running post authentication actions: %v", s.PostAuthenticationActions)
		session.Execute(s.session, run.ParseActions(s.PostAuthenticationActions...)...)
	}
}

func (s *Session) RunWith(token string, fn func()) {
	s.Run(token)

	// after authenticated, run fn, wiggle mouse ...
	fn()
}

func (s *Session) useLightVersion() {
	log.Info().Msgf("useLightVersion: %v", s.UseLightVersion)

	if !s.UseLightVersion {
		return
	}

	if chromedpexecutor.Exists(s.session, time.Duration(time.Second*5), "protocolhandler-welcome-useLightVersionLink", chromedp.ByID) {
		log.Info().Msg("switching to light version")
		session.Execute(s.session,
			chromedp.Click("protocolhandler-welcome-useLightVersionLink", chromedp.ByID),
		)
	}
}
