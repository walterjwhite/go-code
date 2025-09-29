package action

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"time"
)

type key int

const (
	ContextKey key = iota
)

func Execute(ctx context.Context, actions ...chromedp.Action) error {
	log.Debug().Msgf("running [%v] - %v", ctx, actions)

	now := time.Now().Unix()
	ScreenshotIfDebug(ctx, "/tmp/chromedp-%v-0.png", now)
	defer ScreenshotIfDebug(ctx, "/tmp/chromedp-%v-1.png", now)

	d := ctx.Value(ContextKey)
	if d == nil {
		return chromedp.Run(ctx, actions...)
	}

	delayer, ok := d.(delay.Delayer)
	if !ok {
		return errors.New("not a delayer")
	}

	return ExecuteWithDelay(ctx, delayer, actions...)
}

func ExecuteWithDelay(ctx context.Context, delay delay.Delayer, actions ...chromedp.Action) error {
	for i, action := range actions {
		log.Debug().Msgf("running %v with delay", action)

		err := chromedp.Run(ctx, action)
		if err != nil {
			return err
		}

		if i < len(actions)-1 {
			delay.Delay()
		}
	}

	return nil
}

func TryExecute(ctx context.Context, attemptCount int, actions ...chromedp.Action) {
	log.Debug().Msgf("running [%v] - %v: %v", ctx, actions, attemptCount)
	for i := 0; i < attemptCount; i++ {
		err := chromedp.Run(ctx, actions...)
		if err == nil {
			return
		}
	}

	logging.Panic(fmt.Errorf("failed after %v attempts", attemptCount))
}
