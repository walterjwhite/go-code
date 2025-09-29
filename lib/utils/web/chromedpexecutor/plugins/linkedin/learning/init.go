package learning

import (
	"context"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

func (s *Session) Init(ctx context.Context) {
	s.ctx, s.cancel = provider.New(s.Conf, ctx)
}
