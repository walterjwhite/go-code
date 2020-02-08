package pnc

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	logoutButton = "//*[@id=\"topLinks\"]/ul/li[3]/a"

	usernameField = "//*[@id=\"userId\"]"
	passwordField = "//*[@id=\"passwordInputField\"]"
	//loginButton   = "//*[@id=\"olb-btn\"]"
)

// this works partially
// need to reimplement detecting secret questions and then auto-supplying the answers
func (s *PNCSession) Authenticate(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Submit(passwordField),
	)

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
			Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	)
}

func (s *PNCSession) Logout() {
	defer s.chromedpsession.Cancel()

	s.chromedpsession.Execute(chromedp.Click(logoutButton))
}
