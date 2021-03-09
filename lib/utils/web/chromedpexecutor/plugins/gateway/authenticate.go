package gateway

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
)

const (
	menuButtonXpath   = "//*[@id=\"userMenuBtn\"]/div"
	logoffButtonXpath = "//*[@id=\"menuLogOffBtn\"]"
)

// authenticate and nothing more
func (s *Session) Authenticate(ctx context.Context) {
	if len(s.Token) != 6 {
		logging.Panic(fmt.Errorf("Please enter the 6-digit token: %v", s.Token))
	}

	s.ChromeDPSession = chromedpexecutor.New(ctx)

	// no need to wait
	s.ChromeDPSession.Waiter = nil

	s.ChromeDPSession.Execute(chromedp.Navigate(s.Endpoint.Uri))

	log.Debug().Msgf("username: %v", s.Credentials.Username)
	log.Debug().Msgf("domain: %v", s.Credentials.Domain)
	log.Debug().Msgf("password: %v", s.Credentials.Password)
	log.Debug().Msgf("pin/token: %v", s.getToken())

	s.ChromeDPSession.Execute(
		chromedp.SendKeys(s.Endpoint.UsernameXPath, s.Credentials.Domain+"\\"+s.Credentials.Username),
		chromedp.SendKeys(s.Endpoint.PasswordXPath, s.Credentials.Password),
		chromedp.SendKeys(s.Endpoint.TokenXPath, s.getToken()),
		//		chromedp.Click(s.Endpoint.LoginButtonXPath),
		//chromedp.Submit(s.Endpoint.TokenXPath),
		chromedp.KeyEvent(kb.Enter),
	)
}

func (s *Session) getToken() string {
	return s.Credentials.Pin + s.Token
}

// TODO: configure this
func (s *Session) Logout() {
	s.ChromeDPSession.Execute(
		chromedp.Click(menuButtonXpath),
		chromedp.Click(logoffButtonXpath),
	)
}

// TODO: configure this
func (s *Session) isAuthenticated() bool {
	return s.ChromeDPSession.Exists(menuButtonXpath)
}
