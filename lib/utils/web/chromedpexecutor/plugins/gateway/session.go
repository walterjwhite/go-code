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
	useLightVersionPromptId = "protocolhandler-welcome-useLightVersionLink"
)

func (s *Session) Run(token string) bool {
	log.Info().Msgf("running with: %v", token)
	validateToken(token)


	s.Authenticate(token)
	if log.Debug().Enabled() {
		chromedpexecutor.FullScreenshot(s.session.Context(), "/tmp/0.gateway-authenticate.png")
	}

	if !s.IsAuthenticated() {
		return false
	}

	if log.Debug().Enabled() {
		chromedpexecutor.FullScreenshot(s.session.Context(), "/tmp/1.gateway-authenticated.png")
	}

	s.useLightVersion()
	s.runPostAuthenticationActions()

	log.Info().Msg("returning from run")
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

func (s *Session) useLightVersion() {
	log.Info().Msgf("useLightVersion: %v", s.UseLightVersion)

	if !s.UseLightVersion {
		return
	}

	if chromedpexecutor.Exists(s.session, time.Duration(time.Second*5), useLightVersionPromptId, chromedp.ByID) {
		log.Info().Msg("switching to light version")
		session.Execute(s.session,
			chromedp.Click(useLightVersionPromptId, chromedp.ByID),
		)
	}
}
