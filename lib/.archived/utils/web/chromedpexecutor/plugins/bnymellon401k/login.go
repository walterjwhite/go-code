package bnymellon401k

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

const (
	// url           = "https://my.voya.com/voyassoui/index.html?domain=bnymellon.voya.com"
	url           = "https://bnymellon.voya.com"
	usernameField = "//*[@id=\"emailOrUsername\"]"
	passwordField = "//*[@id=\"password\"]"

	loginButton = "//*[@id=\"doc-main-inner\"]/compose/div/div[1]/login-block/form/div/voya-button/span"
)

func (s *Session) Login(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// we must wait, the system prevents / dislikes bots
	//s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		// chromedp.Submit(passwordField),
		chromedp.Click(loginButton),
	)
}
