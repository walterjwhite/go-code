package craigslist

import (
	"github.com/chromedp/chromedp"
)

type OwnerType int

const (
	Owner OwnerType = iota
)

func (p *CraigslistPost) doForSaleBy() []chromedp.Action {
	// for sale by owner
	return []chromedp.Action{chromedp.Click("/html/body/article/section/form/ul/li[6]/label/span[1]/input")}
}
