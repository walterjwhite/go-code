package static

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

type Conf struct {
	InitialActionDelay time.Duration
	TimeBetweenActions time.Duration
	Iterations         int
}

func (c *Conf) Handle(ctx context.Context) error {
	log.Info().Msgf("handling prompt - static")

	time.Sleep(c.InitialActionDelay)

	log.Info().Msg("clicking @ 100,100")
	err := chromedp.Run(ctx,
		chromedp.MouseClickXY(100, 100))
	logging.Warn(err, false, "static.Handle")
	if err != nil {
		return err
	}

	log.Info().Msg("clicked @ 100,100")

	for iteration := 0; iteration < c.Iterations; iteration++ {
		log.Info().Msgf("hitting enter - %d", iteration)
		err = chromedp.Run(ctx,
			chromedp.KeyEvent(kb.Enter),
			chromedp.Sleep(c.TimeBetweenActions))

		logging.Warn(err, false, "static.Handle.iteration")
		if err != nil {
			return err
		}
	}

	log.Info().Msg("done - static")
}
