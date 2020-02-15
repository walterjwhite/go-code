package gateway

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"github.com/walterjwhite/go-application/libraries/logging"
)

// authenticate and nothing more
func (s *Session) Authenticate(ctx context.Context) {
	if len(s.Token) != 6 {
		logging.Panic(fmt.Errorf("Please enter the 6-digit token: %v", s.Token))
	}

	s.chromedpsession = chromedpexecutor.LaunchRemoteBrowser(ctx)

	// no need to wait
	s.chromedpsession.Waiter = nil

	s.chromedpsession.Execute(chromedp.Navigate(s.Endpoint.Uri))

	log.Debug().Msgf("pin: %v%v", s.Credentials.Pin, s.Token)

	s.chromedpsession.Execute(
		chromedp.SendKeys(s.Endpoint.UsernameXPath, s.Credentials.Domain+"\\"+s.Credentials.Username),
		chromedp.SendKeys(s.Endpoint.PasswordXPath, s.Credentials.Password),
		chromedp.SendKeys(s.Endpoint.TokenXPath, s.getToken()),
		chromedp.Click(s.Endpoint.LoginButtonXPath),
	)
}

func (s *Session) getToken() string {
	return s.Credentials.Pin + s.Token
}

// TODO: configure this
func (s *Session) Logout() {
	s.chromedpsession.Execute(
		chromedp.Click("//*[@id=\"userMenuBtn\"]/p"),
		chromedp.Click("//*[@id=\"menuLogOffBtn\"]"),
	)
}

// TODO: configure this
func (s *Session) isAuthenticated() bool {
	return s.chromedpsession.Exists("//*[@id=\"userMenuBtn\"]/p")
}
