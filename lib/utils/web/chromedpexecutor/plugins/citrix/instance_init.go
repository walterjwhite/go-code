package citrix

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"sync/atomic"
	"time"
)

func (i *Instance) init() {
	if i.isInitialized() {
		log.Debug().Msgf("%v - Instance.init - already initialized", i)

		logging.Warn(i.unlock(), false, "Instance.init - error unlocking")
		return
	}

	i.active = &atomic.Bool{}

	i.launch()

	i.acceptTerms()
	i.waitForSessionReady()

	closePermissionPrompts(i.ctx)

	if !i.session.Conf.Headless {
		logging.Warn(action.AttachMousePositionListener(i.ctx), false, "AttachMousePositionListener")
	}

	i.initializeWorker()
	i.makeFullScreen()
	i.actions()
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

	ctx, cancel := context.WithTimeout(i.ctx, 1*time.Second)
	defer cancel()

	logging.Warn(chromedp.Run(ctx, chromedp.MouseClickXY(100, 100), chromedp.KeyEvent(kb.Enter)), false, "Instance.acceptTerms - error accepting terms")

	log.Info().Msgf("%v - Instance.acceptTerms - end", i)
}

func (i *Instance) waitForSessionReady() {
	citrixSessionReadyCtx, citrixSessionReadyCancel := context.WithTimeout(i.ctx, citrixSessionInitializationTimeout)
	defer citrixSessionReadyCancel()

	err := waitForCitrixSessionReady(citrixSessionReadyCtx)
	logging.Warn(err, false, "Instance.waitForSessionReady - error waiting for session ready")
}

func (i *Instance) makeFullScreen() {
	if !i.FullScreen {
		log.Info().Msgf("%v - Instance.makeFullScreen - fullscreen not requested", i)
		return
	}

	successful := action.Fullscreen(i.ctx)
	log.Info().Msgf("%v - Instance.makeFullScreen - fullscreen: %v", i, successful)
}
