package session

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type ChromeDPSession interface {
	Context() context.Context
	Cancel()
}

// type ChromeDPSession struct {
// 	context context.Context
// 	cancel  context.CancelFunc

// 	waiter sleep.Waiter

// 	limit *time.Duration
// }

// func (s *ChromeDPSession) Cancel() {
// 	s.cancel()
// }

func Execute(s ChromeDPSession, actions ...chromedp.Action) {
	logging.Panic(chromedp.Run(s.Context(), actions...))
}
