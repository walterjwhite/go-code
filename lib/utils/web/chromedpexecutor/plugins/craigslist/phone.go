package craigslist

import (
	"github.com/chromedp/chromedp"
)

func (p *CraigslistPost) doPhone() []chromedp.Action {
	var actions = make([]chromedp.Action, 0)

	if len(p.PhoneNumber) == 0 {
		return actions
	}

	// show phone number
	actions = append(actions, chromedp.Click("//*[@id=\"new-edit\"]/div/fieldset[2]/div/fieldset/div/label/input"))

	// phone calls okay
	if p.ReceiveCalls {
		actions = append(actions, chromedp.Click("//*[@id=\"new-edit\"]/div/fieldset[2]/div/fieldset/div/div[1]/label[1]/div/span"))
	}

	// text / sms okay
	if p.ReceiveTexts {
		chromedp.Click("//*[@id=\"new-edit\"]/div/fieldset[2]/div/fieldset/div/div[1]/label[2]/div/span")
	}

	// phone number
	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/fieldset[2]/div/fieldset/div/div[2]/label[1]/label/input", p.PhoneNumber))

	// contact name
	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/fieldset[2]/div/fieldset/div/div[2]/label[3]/label/input", p.PhoneContactName))

	// continue
	actions = append(actions, chromedp.Click("//*[@id=\"new-edit\"]/div/div[3]/div/button"))

	return actions
}
