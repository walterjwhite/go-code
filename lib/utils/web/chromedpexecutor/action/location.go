package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"time"
)

func Location(pctx context.Context) string {
	log.Debug().Msg("location - start")
	ctx, cancel := context.WithTimeout(pctx, 250*time.Millisecond)
	defer cancel()

	var currentUrl string
	logging.Warn(chromedp.Run(ctx, chromedp.Location(&currentUrl)), "Location()")

	log.Debug().Msgf("current url: %s", currentUrl)
	return currentUrl
}
