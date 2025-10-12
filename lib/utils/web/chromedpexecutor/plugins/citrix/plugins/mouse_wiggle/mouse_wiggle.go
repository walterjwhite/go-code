package mouse_wiggle

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

func (c *Conf) Name() string {
	return "mouse wiggle"
}

func (c *Conf) Work(ctx context.Context, headless bool) {
	log.Debug().Msgf("Conf.Work - mouse wiggle is enabled: %s", action.Location(ctx))
	defer log.Debug().Msg("Conf.Work - mouse wiggle - done")
	defer action.ScreenshotIfDebug(ctx, "/tmp/gateway-mouse-wiggle-%v.png", time.Now().Unix())

	for _, coordinates := range c.Points {
		err, wasMoved := action.WasMouseMoved(ctx)
		logging.Warn(err, false, "Conf.Work - error checking if mouse was moved")
		if wasMoved {
			log.Warn().Msg("Conf.Work - mouse was moved")
			return
		}

		c.moveMouse(ctx, coordinates.X, coordinates.Y)
	}
}

func (c *Conf) moveMouse(ctx context.Context, x, y float64) {
	err := action.MoveMouse(ctx, x, y)
	logging.Warn(err, false, "Conf.moveMouse - error moving mouse")

	time.Sleep(c.TimeBetweenActions)
}

func (c *Conf) Cleanup() {
}
