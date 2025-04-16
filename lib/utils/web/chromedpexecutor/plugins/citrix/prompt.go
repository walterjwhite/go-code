package citrix

import (
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (i *Instance) handlePrompt() {
	log.Info().Msgf("handling prompt - %d", i.Index)

	logging.Panic(chromedp.Run(i.ctx,
		chromedp.MouseEvent(input.MouseMoved, 100, 100)))

	for iteration := 0; iteration < 3; iteration++ {
		log.Info().Msgf("hitting enter - %d:%d", i.Index, iteration)
		logging.Panic(chromedp.Run(i.ctx,
			chromedp.KeyEvent(kb.Enter),
			chromedp.Sleep(*i.TimeBetweenActions)))
	}

	log.Info().Msgf("handled prompt - %d", i.Index)
	i.captureScreenshot("/tmp/2.gateway-prompt-%d.png")
}
