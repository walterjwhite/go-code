package citrix

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

func (s *Session) Init(ctx context.Context) {
	if !s.Worker.WillRun() {
		log.Warn().Msg("will not run")
		return
	}

	s.ctx, s.cancel = provider.New(s.Conf, ctx)
}

func (s *Session) Runnable() bool {
	return s.ctx != nil
}
