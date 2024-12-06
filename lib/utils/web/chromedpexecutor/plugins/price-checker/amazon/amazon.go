package amazon

import (
	"strings"

	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

type Amazon struct {
	Url     *string
	Session session.ChromeDPSession
}

func New(s session.ChromeDPSession, url *string) *Amazon {
	return &Amazon{Url: url, Session: s}
}

func (a *Amazon) IsInStock() bool {
	var availability string
	return strings.Contains(availability, "In Stock")
}

func (a *Amazon) GetPrice() string {
	var price string

	if len(price) == 0 {
	}

	return strings.TrimSpace(price)
}
