package headless

import (
	"context"

	"github.com/chromedp/chromedp"
)

type HeadlessChromeDPSession struct {
	ctx               context.Context
	ctxCancelFunction context.CancelFunc
}

func New(ctx context.Context) *HeadlessChromeDPSession {
	ctx1, cancel := chromedp.NewContext(ctx)
	return &HeadlessChromeDPSession{ctx: ctx1, ctxCancelFunction: cancel}
}

func (s *HeadlessChromeDPSession) Context() context.Context {
	return s.ctx
}

func (s *HeadlessChromeDPSession) Cancel() {
	s.ctxCancelFunction()
}
