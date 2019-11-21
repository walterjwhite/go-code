package craigslist

import (
	"github.com/chromedp/chromedp"
)

func (p *CraigslistPost) doMedia() []chromedp.Action {
	var actions []chromedp.Action

	for _, image := range p.Media {
		actions = append(actions, chromedp.SendKeys("#plupload", image))
	}

	// add images
	//chromedp.Click("//*[@id=\"plupload\"]"),
	//chromedp.SetFileInputFiles("#plupload", p.media),
	//chromedp.SendKeys("#plupload", p.Media),

	// done with images
	actions = append(actions, chromedp.Click("/html/body/article/section/form/button"))

	return actions
}
