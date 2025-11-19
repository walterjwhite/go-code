package citrix

import (
	"context"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"time"
)

const (
	launchTimeout       = 1 * time.Minute
	launchRetryAttempts = 5
	launchRetryDelay    = 5 * time.Second

	instanceLaunchXPath = "//*[@id=\"home-screen\"]/div[2]/section[5]/div[5]/div/ul/li[%d]/a[1]"

	tabType = "page"
)

func (i *Instance) launch() {
	log.Debug().Msgf("%v - Instance.launch - start", i)

	if IsContextExpired(i.session.ctx) {
		log.Warn().Msgf("%v - Instance.launch - session appears to have expired", i)
		return
	}

	err := retry.Do(
		func() error {
			return i.tryLaunch()
		},
		retry.Attempts(launchRetryAttempts),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return retry.BackOffDelay(n, err, config)
		}),
		retry.Delay(launchRetryDelay),
	)

	logging.Warn(err, "Instance.launch - error launching")
	if err != nil {
		return
	}

	action.OnTabClosed(i.ctx, i.cleanup)

	log.Debug().Msgf("%v - Instance.launch - end", i)
}

func (i *Instance) tryLaunch() error {
	log.Debug().Msgf("%v - Instance.tryLaunch - launching instance: %d @ [%s]", i, i.Index, action.Location(i.session.ctx))
	targetElementXpath := fmt.Sprintf(instanceLaunchXPath, i.Index)
	targetIDChannel := chromedp.WaitNewTarget(i.session.ctx, matchTabWithNonEmptyURL)

	ctx, cancel := context.WithTimeout(i.session.ctx, launchTimeout)
	defer cancel()

	log.Debug().Msgf("%v - Instance.tryLaunch - clicking: %s", i, targetElementXpath)
	err := chromedp.Run(ctx, chromedp.Click(targetElementXpath))
	if err != nil {
		return err
	}

	log.Debug().Msgf("%v - Instance.tryLaunch - clicked", i)

	select {
	case targetID := <-targetIDChannel:
		tabCtx, tabCancel := chromedp.NewContext(i.session.ctx, chromedp.WithTargetID(targetID))
		err = chromedp.Run(tabCtx)
		logging.Warn(err, "Instance.tryLaunch - error marking context for tab")
		if err != nil {
			tabCancel()
			return err
		}

		log.Debug().Msgf("%v - Instance.tryLaunch - new instance", i)

		logging.Warn(i.waitForInitialization(tabCtx), "Instance.tryLaunch - error waiting for initialization")

		log.Info().Msgf("%v - Instance.tryLaunch - instance successfully initialized", i)

		i.ctx = tabCtx
		i.cancel = tabCancel

		log.Debug().Msgf("%v - Instance.tryLaunch - end - success", i)

		return nil
	case <-ctx.Done():
		log.Debug().Msgf("%v - Instance.tryLaunch - launch context done", i)

		return ctx.Err()
	}
}

func matchTabWithNonEmptyURL(info *target.Info) bool {
	return info.Type == tabType && info.URL != ""
}
