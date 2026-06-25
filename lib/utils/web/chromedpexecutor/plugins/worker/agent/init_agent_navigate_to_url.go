package agent

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/until"
	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"time"
)

func (c *Conf) navigateToUrl(pctx context.Context) error {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	i := graphical.ImageMatch{Ctx: ctx, Image: edgeIconImage, Controller: c.WindowsConf.Controller}
	err := until.Until(ctx, 250*time.Millisecond, i.Matches)
	if err != nil {
		log.Warn().Msgf("agent.navigateToUrl.isPresent.error - %v", err)
		return err
	}

	log.Info().Msgf("agent.navigateToUrl.imageMatched: %s", c.Url)
	return action.Execute(ctx,
		chromedp.MouseClickXY(float64(i.Match.Rect.Max.X)+150,
			float64(i.Match.Rect.Min.Y)+46),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.KeyEvent(c.Url),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.KeyEvent(kb.Enter))
}
