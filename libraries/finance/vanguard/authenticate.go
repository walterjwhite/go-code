package vanguard

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	logoutButton = "//*[@id=\"globalNavUtilityBar\"]/div/div/ul/li[5]/a/span"

	personalInvestors = "/html/body/app-root/app-home-page/homepage-get-started/section/div/div/aside/div/div/ul/app-triage/li[1]/div/span[1]"

	usernameField = "//*[@id=\"username\"]"
	passwordField = "//*[@id=\"password\"]"
	//loginButton   = "//*[@id=\"olb-btn\"]"
)

// this works partially
// need to reimplement detecting secret questions and then auto-supplying the answers
func (s *VanguardSession) Authenticate(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.LaunchRemoteBrowser(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.Click(personalInvestors),
	)

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(usernameField),
			Limit: 5 * time.Second, IsException: true, Message: "Login Form is visible"},
	)

	s.chromedpsession.Execute(
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Submit(passwordField),
	)

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
			Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	)
}

func (s *VanguardSession) Logout() {
	defer s.chromedpsession.Cancel()

	s.chromedpsession.Execute(chromedp.Click(logoutButton))
}
