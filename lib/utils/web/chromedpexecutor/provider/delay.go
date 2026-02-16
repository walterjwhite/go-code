package provider

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

func (c *Conf) withDelayer(pctx context.Context) context.Context {
	if c.Delay <= 0 {
		log.Debug().Msg("no delay configured")
		return pctx
	}

	if c.DelayType == delay.Fixed {
		log.Debug().Msgf("using fixed delay: %v", c.Delay)
		return context.WithValue(pctx, action.ContextKey, delay.New(c.Delay))
	}

	log.Debug().Msgf("using random delay: %v", c.Delay)
	return context.WithValue(pctx, action.ContextKey, delay.NewRandom(c.Delay))
}
