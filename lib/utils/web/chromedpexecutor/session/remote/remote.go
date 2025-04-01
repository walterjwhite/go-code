package remote

import (
	"context"
	"flag"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type RemoteChromeDPSession struct {
	allocatorContext        context.Context
	allocatorCancelFunction context.CancelFunc

	ctx               context.Context
	ctxCancelFunction context.CancelFunc
}

var (
	devToolsWsUrlFlag = flag.String("u", "", "Dev Tools WS URL")
)

func New(ctx context.Context) *RemoteChromeDPSession {
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

	return newRemoteInstance(ctx)
}

func (s *RemoteChromeDPSession) Context() context.Context {
	return s.ctx
}

func (s *RemoteChromeDPSession) Cancel() {
	defer s.allocatorCancelFunction()

	s.ctxCancelFunction()
}

func newRemoteInstance(ctx context.Context) *RemoteChromeDPSession {
	allocatorContext, allocatorCancelFunction := chromedp.NewRemoteAllocator(ctx, *devToolsWsUrlFlag)
	ctx, ctxCancelFunction := chromedp.NewContext(allocatorContext)
	return &RemoteChromeDPSession{allocatorContext: allocatorContext, allocatorCancelFunction: allocatorCancelFunction, ctx: ctx, ctxCancelFunction: ctxCancelFunction}
}
