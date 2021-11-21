package virginpulse

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	//"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"time"
)

// authenticate and nothing more
func (s *Session) Authenticate(ctx context.Context) {
	s.ChromeDPSession = chromedpexecutor.LaunchRemoteBrowser(ctx)

	// no need to wait
	s.ChromeDPSession.Waiter = nil

	s.ChromeDPSession.Execute(chromedp.Navigate(s.Uri))

	log.Info().Msgf("emailAddress: %v", s.Credentials.EmailAddress)

	if s.ByPassAuthentication {
		return
	}

	s.ChromeDPSession.Execute(
		//chromedp.KeyEvent(kb.Control)
		//chromedp.SendKeys(s.UsernameXpath, kb.Delete),
		// do not remember username
		//chromedp.Click(//*[@id="kc-form-wrapper"]/div/form/div/div/div[2]/div[6]/div/label"),
		chromedp.WaitVisible(s.UsernameXpath),
		chromedp.Clear(s.UsernameXpath),
		chromedp.SetValue(s.UsernameXpath, s.Credentials.EmailAddress),
		chromedp.SetValue(s.PasswordXpath, s.Credentials.Password),
		chromedp.Click(s.LoginButtonXpath),
	)

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	if !s.IsAuthenticated() {
		logging.Panic(fmt.Errorf("Should be authenticated, but not"))
	}
}

// TODO: configure this
func (s *Session) Logout() {
	s.ChromeDPSession.Execute(
		chromedp.Click(s.MenuXpath),
		chromedp.Click(s.LogoffButtonXpath),
	)
}

// TODO: configure this
func (s *Session) IsAuthenticated() bool {
	return s.ChromeDPSession.Exists(s.MenuXpath)
}
