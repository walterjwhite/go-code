package headless

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type HeadlessChromeDPSession struct {
	ctx               context.Context
	ctxCancelFunction context.CancelFunc
}

func New(ctx context.Context) *HeadlessChromeDPSession {
	ctx1, cancel := chromedp.NewContext(ctx)

	logging.Panic(chromedp.Run(ctx1, chromedp.EmulateViewport(1920, 1080)))

	return &HeadlessChromeDPSession{ctx: ctx1, ctxCancelFunction: cancel}
}

func (s *HeadlessChromeDPSession) Context() context.Context {
	return s.ctx
}

func (s *HeadlessChromeDPSession) Cancel() {
	s.ctxCancelFunction()
}
