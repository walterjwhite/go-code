package clearviewfcu

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
	// "time"
)

const (
	url = "https://www.clearviewfcu.org"

	usernameField = "//*[@id=\"userid\"]"
	passwordField = "//*[@id=\"password\"]"

	loginButton = "//*[@id=\"LoginButton\"]"
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
		chromedp.Click(loginButton),
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Submit(passwordField),
	)

	// s.chromedpsession.ExecuteTimeLimited(
	// 	chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
	// 		Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	// )
}
