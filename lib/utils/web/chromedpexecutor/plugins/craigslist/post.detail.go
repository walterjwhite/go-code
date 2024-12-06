package craigslist

import (
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
)

func (p *CraigslistPost) doPostDetails() []chromedp.Action {
	var actions []chromedp.Action






	actions = append(actions, p.doScript()...)


	return actions
}

func (p *CraigslistPost) doScript() []chromedp.Action {
	if p.Script == nil {
		return nil
	}

	return run.ParseActions(p.Script...)
}
