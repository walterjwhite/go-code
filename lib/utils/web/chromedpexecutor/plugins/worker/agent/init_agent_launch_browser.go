package agent

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"time"
)

func (c *Conf) launchBrowser(pctx context.Context) error {
	if c.WindowsConf != nil {
		log.Info().Msg("agent.Init.launchBrowser.WindowsConf != nil")
		err := c.waitForWindowsToLoad(pctx)
		if err != nil {
			return err
		}

		log.Info().Msg("agent.Init.launchBrowser.WindowsConf - done")
	}

	ctx, cancel := context.WithTimeout(pctx, 15*time.Second)
	defer cancel()

	log.Info().Msgf("agent.launchBrowser: %s", c.BrowserName)
	return action.Execute(ctx,
		chromedp.KeyEvent(kb.Meta),
		chromedp.Sleep(1*time.Second),
		chromedp.KeyEvent(c.BrowserName),
		chromedp.Sleep(1*time.Second),
		chromedp.KeyEvent(kb.Enter),
		chromedp.Sleep(5*time.Second))
}
