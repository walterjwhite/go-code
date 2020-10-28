package chromedpexecutor

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"time"
)

type TimeLimitedChromeAction struct {
	Action      chromedp.Action
	Limit       time.Duration
	IsException bool
	Message     string
}

func (s *ChromeDPSession) ExecuteTimeLimited(actions ...TimeLimitedChromeAction) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)

		ctx, cancel := context.WithTimeout(s.Context, action.Limit)
		defer cancel()

		logging.Warn(chromedp.Run(ctx, action.Action), action.IsException, action.Message)

		if i < (len(actions) - 1) {
			s.Waiter.Wait()
		}
	}
}
