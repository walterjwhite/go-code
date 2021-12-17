package price_checker

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

type PriceChecker interface {
	// select options ...

	IsInStock() bool
	GetPrice() string
}

type Session struct {
	Session      *chromedpexecutor.ChromeDPSession
	Url          *string
	PriceChecker PriceChecker
}

func New(ctx context.Context, url *string) *Session {
	s := chromedpexecutor.New(ctx)
	s.Waiter = nil

	return &Session{Session: s, Url: url}
}

func (s *Session) Fetch() string {
	s.Session.Execute(chromedp.Navigate(*s.Url))

	if !s.PriceChecker.IsInStock() {
		logging.Panic(fmt.Errorf("Item is not in stock: %s", *s.Url))
	}

	return s.PriceChecker.GetPrice()
}
