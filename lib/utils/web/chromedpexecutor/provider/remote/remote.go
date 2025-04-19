package remote

import (
	"context"
	"flag"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	devToolsWsUrlFlag = flag.String("u", "", "Dev Tools WS URL")

	allocatorContext context.Context
	allocatorCancel  context.CancelFunc
)

func New(ctx context.Context) (context.Context, context.CancelFunc) {
	initAllocator(ctx)
	ictx, _ := chromedp.NewContext(allocatorContext)
	return ictx, allocatorCancel
}

func initAllocator(ctx context.Context) {
	if allocatorContext != nil {
		return
	}

	if len(*devToolsWsUrlFlag) == 0 {
		if len(*devToolsWsFileFlag) > 0 {
			getURLFromFile()
		}

		if len(*devToolsWsUrlFlag) == 0 {
			log.Info().Msg("launching new instance")
			launchRemoteBrowser(ctx)

			getURLFromFile()

			if len(*devToolsWsUrlFlag) == 0 {
				logging.Panic(fmt.Errorf("unable to determine dev tools url"))
			}
		}
	}

	allocatorContext, allocatorCancel = chromedp.NewRemoteAllocator(ctx, *devToolsWsUrlFlag)
}
