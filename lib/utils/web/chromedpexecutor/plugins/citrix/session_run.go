package citrix

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/chromedp/cdproto/browser"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

func (s *Session) Run(token string) error {
	log.Debug().Msg("Session.Run - start")

	defer s.cancel()

	var javascriptListenCancel context.CancelFunc
	if log.Debug().Enabled() {
		javascriptListenCancel = action.CaptureJavascript(s.ctx)
	}

	err := s.authenticate(token)
	if err != nil {
		return err
	}

	s.useLightVersion()

	action.Grant(s.ctx, []browser.PermissionType{"windowManagement"})

	if javascriptListenCancel != nil {
		javascriptListenCancel()
	}

	go s.keepAlive()
	go s.onDone()


	s.Worker.WithWorker(s)
	s.Worker.Run(s.ctx)

	log.Debug().Msg("Session.Run - end")
	return nil
}

func (s *Session) onDone() {
	<-s.ctx.Done()
	s.cleanup()
}

func (s *Session) cleanup() {
	log.Info().Msg("Session.cleanup - cleaning up")
	logging.Warn(s.GoogleProvider.PublishStatus("cleanup session", true), "session.cleanup")
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

		logging.Warn(s.Instances[i].lock(), "Session.lockWorkers - error locking worker")
	}
}
