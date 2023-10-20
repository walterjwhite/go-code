package jensonusa

import (
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

type Jenson struct {
	Url     *string
	session session.ChromeDPSession
}

func New(s session.ChromeDPSession, url *string) *Jenson {
	return &Jenson{Url: url, session: s}
}

func (j *Jenson) IsInStock() bool {
	err := false
	defer func() {
		if r := recover(); r != nil {
			err = true
		}
	}()

	var availability string
	session.Execute(j.session, chromedp.Text("txtAddToCart", &availability, chromedp.NodeVisible, chromedp.ByID))
	return !err
}

func (j *Jenson) GetPrice() string {
	var price string
	session.Execute(j.session, chromedp.Text("price", &price, chromedp.NodeVisible, chromedp.ByID))
	return strings.TrimSpace(price)
}
