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
	session.Execute(a.Session, chromedp.Text("//*[@id=\"availability\"]/span", &availability, chromedp.NodeVisible, chromedp.BySearch))
	return strings.Contains(availability, "In Stock")
}

func (a *Amazon) GetPrice() string {
	var price string
	session.Execute(a.Session, chromedp.Text("//*[@id=\"corePrice_desktop\"]/div/table/tbody/tr[2]/td[2]/span[1]/span[2]", &price, chromedp.NodeVisible, chromedp.BySearch))

	if len(price) == 0 {
		session.Execute(a.Session, chromedp.Text("//*[@id=\"corePrice_desktop\"]/div/table/tbody/tr/td[2]/span[1]/span[2]", &price, chromedp.NodeVisible, chromedp.BySearch))
	}

	return strings.TrimSpace(price)
}
