package gateway

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"context"
	"github.com/chromedp/cdproto/input"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"time"
)

func (s *Session) wiggleMouse(ctx context.Context, instance Instance) {
	if !instance.WiggleMouse {
		log.Warn().Msgf("Mouse wiggle is disabled - %d", instance.Index)
		return
	}

	timeBetweenActions := time.Duration(1 * time.Second)
	log.Warn().Msgf("Mouse wiggle is enabled - %d", instance.Index)

	for {
		s.moveMouse(ctx, 100, 100, instance.Index, timeBetweenActions)
		s.moveMouse(ctx, 200, 100, instance.Index, timeBetweenActions)
		s.moveMouse(ctx, 200, 200, instance.Index, timeBetweenActions)
		s.moveMouse(ctx, 100, 200, instance.Index, timeBetweenActions)

		if log.Debug().Enabled() {
			chromedpexecutor.FullScreenshot(ctx, fmt.Sprintf("/tmp/3.gateway-mouse-wiggle-%d.png", instance.Index))
		}
	}
}

func (s *Session) moveMouse(ctx context.Context, x, y float64, i int, timeBetweenActions time.Duration) {
	log.Info().Msgf("moving mouse to: %f,%f - %d", x, y, i)
	logging.Panic(chromedp.Run(ctx,
		chromedp.MouseEvent(input.MouseMoved, x, y),
		chromedp.Sleep(timeBetweenActions)))
}
