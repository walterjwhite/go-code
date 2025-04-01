package gateway

import (
	"context"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

func handlePrompt(ctx context.Context, instance Instance) {
	log.Info().Msgf("handling prompt - %d", instance.Index)

	logging.Panic(chromedp.Run(ctx,
		chromedp.Sleep(time.Duration(20*time.Second)),
		chromedp.MouseEvent(input.MouseMoved, 100, 100)))

	timeBetweenActions := time.Duration(1 * time.Second)
	for iteration := 0; iteration < 3; iteration++ {
		log.Info().Msgf("hitting enter - %d:%d", instance.Index, iteration)
		logging.Panic(chromedp.Run(ctx,
			chromedp.KeyEvent(kb.Enter),
			chromedp.Sleep(timeBetweenActions)))
	}

	log.Info().Msgf("handled prompt - %d", instance.Index)
}
