package walgreens

import (

	"github.com/chromedp/chromedp"

)

const (
	logoutButton = "//*[@id=\"signOut\"]/strong"
)

func (s *Session) Logout() {
	s.chromedpsession.Execute(
		chromedp.Click(logoutButton),
	)

	// handle popup ...
}