package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

func Execute(ctx context.Context, actions ...chromedp.Action) {
	log.Info().Msgf("running [%v] - %v", ctx, actions)
	logging.Panic(chromedp.Run(ctx, actions...))
}

func ExecuteWithDelay(ctx context.Context, delay delay.Delayer, actions ...chromedp.Action) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)

		logging.Panic(chromedp.Run(ctx, action))

		if i < len(actions)-1 {
			delay.Delay()
		}
	}
}
