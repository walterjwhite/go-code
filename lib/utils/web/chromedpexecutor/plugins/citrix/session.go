package citrix

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/chromedp/cdproto/browser"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"

	"time"
)

const (
	useLightVersionPromptId = "protocolhandler-welcome-useLightVersionLink"
)

func (s *Session) Run(token string) {
	log.Debug().Msg("Session.Run - start")

	defer s.cancel()

	var javascriptListenCancel context.CancelFunc
	if log.Debug().Enabled() {
		javascriptListenCancel = s.captureJavascript()
	}

	s.authenticate(token)
	s.useLightVersion()

	action.Grant(s.ctx, []browser.PermissionType{"windowManagement"})

	s.runPostAuthenticationActions()
	if javascriptListenCancel != nil {
		javascriptListenCancel()
	}

	s.keepAliveTicker = time.NewTicker(*s.Timeout)
	go s.keepAlive()
	go s.onDone()


	s.Worker.WithWorker(s)
	s.Worker.Run(s.ctx)

	log.Debug().Msg("Session.Run - end")
}

func (s *Session) onDone() {
	<-s.ctx.Done()
	s.cleanup()
}

func (s *Session) runPostAuthenticationActions() {
	log.Debug().Msg("Session.runPostAuthenticationActions - start")

	if len(s.PostAuthenticationActions) > 0 {
		log.Info().Msgf("Session.runPostAuthenticationActions - running post authentication actions - delay: %v", *s.PostAuthenticationDelay)
		time.Sleep(*s.PostAuthenticationDelay)

		log.Info().Msgf("Session.runPostAuthenticationActions - running post authentication actions: %v", s.PostAuthenticationActions)
		logging.Warn(action.Execute(s.ctx, run.ParseActions(s.PostAuthenticationActions...)...), false, "session.runPostAuthenticationActions")
	}

	log.Debug().Msg("Session.runPostAuthenticationActions - end")
}

func (s *Session) useLightVersion() {
	log.Info().Msgf("Session.useLightVersion - UseLightVersion: %v", s.UseLightVersion)

	if !s.UseLightVersion {
		return
	}

	select {
	case <-s.ctx.Done():
		log.Debug().Msg("Session.useLightVersion - context done")
	default:
	}

	if action.ExistsById(s.ctx, useLightVersionPromptId) {
		log.Info().Msg("Session.useLightVersion - switching to light version")
		logging.Warn(action.Execute(s.ctx,
			chromedp.Click(useLightVersionPromptId, chromedp.ByID),
		), false, "session.useLightVersion - error selecting use light version")
	}
}

func (s *Session) cleanup() {
	log.Info().Msg("Session.cleanup - cleaning up")
	logging.Warn(s.GoogleProvider.PublishStatus("cleanup session", true), false, "session.cleanup")

	s.keepAliveTicker.Stop()
}

func (s *Session) lockWorkers() {
	for i := range s.Instances {
		if s.Instances[i].active == nil {
			log.Warn().Msg("Session.lockWorkers - instance has not yet been initialized")
			continue
		}

		if s.Instances[i].active.Load() {
			log.Warn().Msgf("Session.lockWorkers - worker is active, not locking: %d", i)
			continue
		}

		log.Warn().Msgf("Session.lockWorkers - locking worker: %d", i)

		logging.Warn(s.Instances[i].lock(), false, "Session.lockWorkers - error locking worker")
	}
}
