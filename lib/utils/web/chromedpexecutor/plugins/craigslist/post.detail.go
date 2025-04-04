package craigslist

import (
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
)

func (p *CraigslistPost) doPostDetails() []chromedp.Action {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys("//*[@id=\"PostingTitle\"]", p.Title))
	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/div[1]/label[2]/label/input", p.Price))

	actions = append(actions, chromedp.SendKeys("//*[@id=\"geographic_area\"]", p.City))

	actions = append(actions, chromedp.SendKeys("//*[@id=\"postal_code\"]", p.PostalCode))

	actions = append(actions, chromedp.SendKeys("//*[@id=\"PostingBody\"]", p.Description))

	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/fieldset[2]/div/div/div[1]/label/label/input", p.EmailAddress))

	actions = append(actions, p.doScript()...)

	actions = append(actions, chromedp.Click("//*[@id=\"new-edit\"]/div/div[3]/div/button"))
	actions = append(actions, chromedp.Click("//*[@id=\"leafletForm\"]/button"))

	return actions
}

func (p *CraigslistPost) doScript() []chromedp.Action {
	if p.Script == nil {
		return nil
	}

	return run.ParseActions(p.Script...)
}
