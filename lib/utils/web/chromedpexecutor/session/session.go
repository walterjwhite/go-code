package session

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
)

type ChromeDPSession interface {
	Context() context.Context
	Cancel()
}





func Execute(s ChromeDPSession, actions ...chromedp.Action) {
	for _, action := range actions {
		log.Info().Msgf("running %v", action)
		logging.Panic(chromedp.Run(s.Context(), action))
	}
}

func ExecuteWithDelay(s ChromeDPSession, delay delay.Delayer, actions ...chromedp.Action) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)

		logging.Panic(chromedp.Run(s.Context(), action))

		if i < len(actions)-1 {
			delay.Delay()
		}
	}
}
