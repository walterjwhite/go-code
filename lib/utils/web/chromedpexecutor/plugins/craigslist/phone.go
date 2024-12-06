package craigslist

import (
	"github.com/chromedp/chromedp"
)

func (p *CraigslistPost) doPhone() []chromedp.Action {
	var actions = make([]chromedp.Action, 0)

	if len(p.PhoneNumber) == 0 {
		return actions
	}


	if p.ReceiveCalls {
	}

	if p.ReceiveTexts {
	}




	return actions
}
