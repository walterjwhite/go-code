package gateway

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (s *Session) Authenticate(ctx context.Context) {
	if len(s.Token) != 6 {
		logging.Panic(errors.New("Please enter the token"))
	}

	s.chromedpsession = chromedpexecutor.New(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	s.chromedpsession.Execute(chromedp.Navigate(s.Endpoint.Uri))

	s.chromedpsession.Execute([]chromedp.Action{
		chromedp.SendKeys(s.Endpoint.UsernameXPath, s.Credentials.Domain+"\\"+s.Credentials.Username),
		chromedp.SendKeys(s.Endpoint.PasswordXPath, s.Credentials.Password),
		chromedp.SendKeys(s.Endpoint.TokenXPath, s.Credentials.Pin+s.Token),
		chromedp.Click(s.Endpoint.LoginButtonXPath),
	}...)

	s.tickle(ctx)
}
