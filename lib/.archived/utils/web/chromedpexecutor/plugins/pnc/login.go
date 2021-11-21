package pnc

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	// "time"
)

const (
	url           = "https://www.pnc.com"
	usernameField = "//*[@id=\"experiencefragment-9df31000f1\"]/div/div/div[3]/div/div/div[2]/form/div/div[1]/div[1]/div[1]/label/input"
	passwordField = "//*[@id=\"experiencefragment-9df31000f1\"]/div/div/div[3]/div/div/div[2]/form/div/div[1]/div[1]/div[3]/label/input"

	loginMenuItem = "//*[@id=\"experiencefragment-9df31000f1\"]/div/div/div[2]/button"
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
		chromedp.Click(loginMenuItem),
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Submit(passwordField),
	)

	// s.chromedpsession.ExecuteTimeLimited(
	// 	chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
	// 		Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	// )
}
