package mouse_driver

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

func (c *Conf) Name() string {
	return "mouse wiggle"
}

func (c *Conf) Init(ctx context.Context, headless bool, contextuals ...interface{}) error {
	return nil
}

func (c *Conf) Work(ctx context.Context, headless bool) {
	log.Debug().Msgf("Conf.Work - mouse wiggle is enabled: %s", action.Location(ctx))
	defer log.Debug().Msg("Conf.Work - mouse wiggle - done")
	defer action.ScreenshotIfDebug(ctx, "/tmp/gateway-mouse-wiggle-%v.png", time.Now().Unix())

	for i, coordinates := range c.Points {
		wasMoved, err := action.WasMouseMoved(ctx)
		logging.Warn(err, fmt.Sprintf("Conf.Work - error checking if mouse was moved [%d]", i))
		if wasMoved {
			log.Warn().Msgf("Conf.Work - mouse was moved in between moveMouse calls [%d]", i)
			return
		}

		c.moveMouse(ctx, coordinates.X, coordinates.Y)
	}
}

func (c *Conf) moveMouse(ctx context.Context, x, y float64) {
	err := action.MoveMouse(ctx, x, y)
	logging.Warn(err, "Conf.moveMouse - error moving mouse")

	time.Sleep(c.TimeBetweenActions)
}

func (c *Conf) Cleanup() {
}
