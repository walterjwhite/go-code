package citrix

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/agent"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/mouse_wiggle"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/noop"

	"sync/atomic"
	"time"
)

func (i *Instance) PostLoad(ctx context.Context) {
	switch i.WorkerType {
	case worker.MouseWiggler:
		i.Worker = &mouse_wiggle.Conf{}
	case worker.Agent:
		i.Worker = &agent.Conf{}
	case worker.NOOP:
		i.Worker = &noop.State{}
	default:
		logging.Panic(errors.New("WorkerType unspecified"))
	}

	application.Load(i.Worker)
}

func (i *Instance) init() {
	if i.isInitialized() {
		log.Debug().Msgf("%v - Instance.init - already initialized", i)

		logging.Warn(i.unlock(), false, "Instance.init - error unlocking")
		return
	}

	i.active = &atomic.Bool{}

	i.launch()

	if i.session.controller == nil {
		i.session.controller = &chromedpexecutor.ChromeDPController{}
	}

	i.WindowsConf.Controller = i.session.controller

	i.acceptTerms()
	i.waitForSessionReady()

	i.closePermissionPrompts()

	if !i.session.Conf.Headless {
		logging.Warn(action.AttachMousePositionListener(i.ctx), false, "AttachMousePositionListener")
	}

	if i.InitializationDelay > 0 {
		log.Info().Msgf("%v - Instance.init - delay - %v", i, i.InitializationDelay)
		time.Sleep(i.InitializationDelay)
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

	visible, err := i.waitForTermsAcceptance()
	if err != nil {
		logging.Warn(err, false, "acceptTerms,waitForTermsAcceptance")
		return
	}

	if !visible {
		log.Warn().Msg("terms acceptance was required, but not found")
		return
	}

	ctx, cancel := context.WithTimeout(i.ctx, 1*time.Second)
	defer cancel()
	logging.Warn(chromedp.Run(ctx, chromedp.MouseClickXY(100, 100), chromedp.KeyEvent(kb.Enter)), false, "Instance.acceptTerms - error accepting terms")

	log.Info().Msgf("%v - Instance.acceptTerms - end", i)
}

func (i *Instance) waitForTermsAcceptance() (bool, error) {
	ctx, cancel := context.WithTimeout(i.ctx, 15*time.Second)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
			visible, err := i.WindowsConf.IsTermsAcceptanceButtonVisible(ctx)
			if err != nil {
				return false, err
			}

			if visible {
				return true, nil
			}
		}
	}
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
