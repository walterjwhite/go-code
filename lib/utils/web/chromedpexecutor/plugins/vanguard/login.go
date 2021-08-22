package vanguard

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
	// "time"
)

const (
	url = "https://investor.vanguard.com/home"

	usernameField = "//*[@id=\"username\"]"
	passwordField = "//*[@id=\"password\"]"
	loginButton   = "//*[@id=\"mainContent\"]/psx-psr2-homepage/section/section/form/div/div/div[3]/vui-button/button"
)

func (s *Session) Login(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
	)

	// s.chromedpsession.ExecuteTimeLimited(
	// 	chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(usernameField),
	// 		Limit: 5 * time.Second, IsException: true, Message: "Login Form is visible"},
	// )

	s.chromedpsession.Execute(
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),

		// this causes the form to be refreshed
		//chromedp.Submit(passwordField),
		// getting a data race here
		chromedp.Click(loginButton),
	)

	// s.chromedpsession.ExecuteTimeLimited(
	// 	chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
	// 		Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	// )
}
