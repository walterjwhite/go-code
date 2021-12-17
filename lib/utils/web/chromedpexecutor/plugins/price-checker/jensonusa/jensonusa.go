package jensonusa

import (
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

type Jenson struct {
	Url     *string
	Session *chromedpexecutor.ChromeDPSession
}

func New(s *chromedpexecutor.ChromeDPSession, url *string) *Jenson {
	return &Jenson{Url: url, Session: s}
}

func (j *Jenson) IsInStock() bool {
	err := false
	defer func() {
		if r := recover(); r != nil {
			err = true
		}
	}()

	var availability string
	j.Session.Execute(chromedp.Text("txtAddToCart", &availability, chromedp.NodeVisible, chromedp.ByID))
	return !err
}

func (j *Jenson) GetPrice() string {
	var price string
	j.Session.Execute(chromedp.Text("price", &price, chromedp.NodeVisible, chromedp.ByID))
	return strings.TrimSpace(price)
}
