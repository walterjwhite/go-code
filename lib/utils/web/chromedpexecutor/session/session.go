package session

import (
	"context"
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
