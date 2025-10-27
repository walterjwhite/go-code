package citrix

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
)

func (i *Instance) run() {
	log.Info().Msgf("%v - Instance.run - start", i)

	defer i.session.waitGroup.Done()
	defer log.Info().Msgf("%v - Instance.run - end", i)

	i.init()
	log.Info().Msgf("%v - Instance.run - running", i)

	err, moved := action.WasMouseMoved(i.ctx)
	if err != nil {
		logging.Warn(err, false, "wasMouseMoved")
		return
	}

	i.active.Store(moved)
	if moved {
		log.Warn().Msgf("%v - Instance.run - mouse was moved, skipping this iteration of Work", i)
	} else {
		ctx, cancel := context.WithTimeout(i.ctx, i.session.Worker.WorkTickInterval/2)
		defer cancel()

		i.Worker.Work(ctx, i.session.Conf.Headless)
	}
}

func (i *Instance) actions() {
	if len(i.Actions) == 0 {
		log.Debug().Msgf("%v - Instance.actions - no actions to run", i)
		return
	}


	log.Info().Msgf("%v - Instance.actions - running actions: %v", i, i.Actions)
	logging.Warn(chromedp.Run(i.ctx, run.ParseActions(i.Actions...)...), false, "Instance.actions - error running actions")
}
