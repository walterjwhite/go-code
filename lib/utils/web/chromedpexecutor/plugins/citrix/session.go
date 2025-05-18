package citrix

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"errors"
	"github.com/chromedp/cdproto/browser"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"

	"time"
)

const (
	useLightVersionPromptId = "protocolhandler-welcome-useLightVersionLink"
)

func (s *Session) Run(token string) {
	defer s.Cancel()

	s.authenticate(token)
	saveScreenshot(s.ctx, "/tmp/%d.gateway-authenticate.png", 0)

	if !s.isAuthenticated() {
		saveScreenshot(s.ctx, "/tmp/%d.gateway-failed-to-authenticate.png", 1)
		logging.Panic(errors.New("failed to authenticate"))
	}

	saveScreenshot(s.ctx, "/tmp/%d.gateway-authenticated.png", 1)

	s.useLightVersion()
	action.Grant(s.ctx, []browser.PermissionType{"windowManagement"})

	s.runPostAuthenticationActions()

	s.keepAliveChannel = time.Tick(*s.Timeout)
	go s.keepAlive()

	s.Worker.Worker = s
	s.Worker.Run()
}

func (s *Session) runPostAuthenticationActions() {
	if len(s.PostAuthenticationActions) > 0 {
		log.Info().Msgf("running post authentication actions - delay: %v", *s.PostAuthenticationDelay)
		time.Sleep(*s.PostAuthenticationDelay)

		log.Info().Msgf("running post authentication actions: %v", s.PostAuthenticationActions)
		action.Execute(s.ctx, run.ParseActions(s.PostAuthenticationActions...)...)
	}
}

func (s *Session) useLightVersion() {
	log.Info().Msgf("useLightVersion: %v", s.UseLightVersion)

	if !s.UseLightVersion {
		return
	}

	if action.Exists(s.ctx, time.Duration(time.Second*5), useLightVersionPromptId, chromedp.ByID) {
		log.Info().Msg("switching to light version")
		action.Execute(s.ctx,
			chromedp.Click(useLightVersionPromptId, chromedp.ByID),
		)
	}
}

func (s *Session) Cancel() {
	log.Warn().Msg("cancelling session")
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	s.cleanup()
	application.Cancel()
}

func (s *Session) cleanup() {
	for index := range s.Instances {
		log.Info().Msgf("on break, cancelling context: %v", s.Instances[index])
		s.Instances[index].cleanup()
	}
}
