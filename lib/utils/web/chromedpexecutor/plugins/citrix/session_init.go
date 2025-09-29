package citrix

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

func (s *Session) Init(ctx context.Context) {
	log.Info().Msg("session.Init(ctx)")

	s.ctx, s.cancel = provider.New(s.Conf, ctx)
}
