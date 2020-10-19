package walgreens

import (

	"github.com/chromedp/chromedp"

)

func (s *Session) Upload(files ...string) {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"menu-photo\"]/a/span"),
		chromedp.Click("//*[@id=\"photoOrg-addPhotos-qmp-btn\"]"),
		chromedp.Click("//*[@id=\"po-fdropdown\"]/li[2]/a/span[1]/span[2]"),
	)

	// handle popup ...
}