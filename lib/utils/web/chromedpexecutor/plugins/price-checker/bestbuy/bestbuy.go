package bestbuy

import (
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

type BestBuy struct {
	Url     *string
	session session.ChromeDPSession
}

func New(s session.ChromeDPSession, url *string) *BestBuy {
	return &BestBuy{Url: url, session: s}
}

func (j *BestBuy) IsInStock() bool {
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

func (j *BestBuy) GetPrice() string {
	var price string
	session.Execute(j.session, chromedp.Text("price", &price, chromedp.NodeVisible, chromedp.ByID))
	return strings.TrimSpace(price)
}
