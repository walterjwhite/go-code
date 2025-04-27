package citrix

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider/headless"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider/remote"
)

func (s *Session) Init(ctx context.Context) {
	if !s.Worker.WillRun() {
		log.Warn().Msg("will not run")
		return
	}

	s.ctx, s.cancel = s.initChromeDP(ctx)
}

func (s *Session) initChromeDP(ctx context.Context) (context.Context, context.CancelFunc) {
	if !s.Headless {
		log.Warn().Msg("New remote session")
		return remote.New(ctx)
	}

	log.Warn().Msg("New headless session")
	return headless.New(ctx)
}

func (s *Session) Runnable() bool {
	return s.ctx != nil
}
