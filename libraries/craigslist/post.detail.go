package craigslist

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/logging"
	"strings"
)

func (p *CraigslistPost) doPostDetails() []chromedp.Action {
	var actions []chromedp.Action

	// title
	actions = append(actions, chromedp.SendKeys("//*[@id=\"PostingTitle\"]", p.Title))
	// price
	// name=price
	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/div[1]/label[2]/label/input", p.Price))

	// city or neighborhood
	actions = append(actions, chromedp.SendKeys("//*[@id=\"GeographicArea\"]", p.City))

	// postal code
	actions = append(actions, chromedp.SendKeys("//*[@id=\"postal_code\"]", p.PostalCode))

	// description
	actions = append(actions, chromedp.SendKeys("//*[@id=\"PostingBody\"]", p.Description))

	// emailAddress
	actions = append(actions, chromedp.SendKeys("//*[@id=\"new-edit\"]/div/fieldset[2]/div/div/div[1]/label/label/input", p.EmailAddress))

	// script
	actions = append(actions, p.doScript()...)

	actions = append(actions, chromedp.Click("//*[@id=\"new-edit\"]/div/div[3]/div/button"))
	actions = append(actions, chromedp.Click("//*[@id=\"leafletForm\"]/button"))

	return actions
}

func (p *CraigslistPost) doScript() []chromedp.Action {
	if p.Script == nil {
		return nil
	}

	actions := make([]chromedp.Action, 0)
	for _, script := range p.Script {
		actions = append(actions, getScript(script))
	}

	return actions
}

func getScript(line string) chromedp.Action {
	arguments := strings.Split(line, ",")

	if arguments[0] == "Click" {
		return chromedp.Click(arguments[1])
	} else if arguments[0] == "SendKeys" {
		return chromedp.SendKeys(arguments[1], arguments[2])
	}

	logging.Panic(fmt.Errorf("Unsupported action: %v", arguments[0]))
	return nil
}
