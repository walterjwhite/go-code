package citrix

import (
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)


func (i *Instance) init() {
	if i.isInitialized() {
		log.Debug().Msgf("%v - Instance.init - already initialized", i)

		logging.Warn(i.unlock(), "Instance.init - error unlocking")
		return
	}

	i.active = &atomic.Bool{}

	i.launch()

	i.setController()

	i.acceptTerms()
	logging.Warn(i.waitForSessionReady(), "Instance.init - error waiting for session ready")

	i.closePermissionPrompts()

	if !i.session.Conf.Headless {
		logging.Warn(action.AttachMousePositionListener(i.ctx), "AttachMousePositionListener")
	}

	i.waitForInitDelay()
	i.makeFullScreen()
	i.actions()
	i.initializeWorker()
}

func (i *Instance) setController() {
	if i.session.controller == nil {
		i.session.controller = &chromedpexecutor.ChromeDPController{}
	}

	i.WindowsConf.Controller = i.session.controller
}

func (i *Instance) isInitialized() bool {
	if i.ctx == nil {
		return false
	}

	select {
	case <-i.ctx.Done():
		return false
	default:
	}

	return true
}

func (i *Instance) acceptTerms() {
	if !i.RequiresTermsAcceptance {
		log.Warn().Msgf("%v - Instance.acceptTerms - does not require terms acceptance", i)
		return
	}

	log.Info().Msgf("%v - Instance.acceptTerms - waiting until ready", i)
	ready, err := i.WindowsConf.WaitReady(i.ctx)
	if ready {
		log.Info().Msgf("%v - acceptTerms.ready", i)
		return
	}

	logging.Warn(err, "acceptTerms.WaitReady")
}

func (i *Instance) makeFullScreen() {
	if !i.FullScreen {
		log.Info().Msgf("%v - Instance.makeFullScreen - fullscreen not requested", i)
		return
	}

	successful := action.Fullscreen(i.ctx)
	log.Info().Msgf("%v - Instance.makeFullScreen - fullscreen: %v", i, successful)
}
