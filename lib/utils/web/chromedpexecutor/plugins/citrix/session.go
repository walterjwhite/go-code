package citrix

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"

	"time"
)

const (
	useLightVersionPromptId = "protocolhandler-welcome-useLightVersionLink"
)

func (s *Session) Run(token string) bool {
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

	s.keepAliveChannel = time.Tick(*s.Timeout)
	go s.keepAlive()

	log.Info().Msg("returning from run")
	return true
}

func (s *Session) scheduleEnd() {
	if s.EndHour == 0 {
		log.Warn().Msg("end hour is unset, running indefinitely")
		return
	}

	waitUntil(s.EndHour)
	application.Cancel()
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
