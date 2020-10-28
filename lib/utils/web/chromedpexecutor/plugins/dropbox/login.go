package discovercard

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
	"time"
)

const (
	url = "https://www.dropbox.com/login"

	usernameField = "/html/body/div[12]/div[1]/div[2]/div/div/div[1]/div[2]/div/div/div/form/div[1]/div[1]/div[2]/input"
	passwordField = "/html/body/div[12]/div[1]/div[2]/div/div/div[1]/div[2]/div/div/div/form/div[1]/div[2]/div[2]/input"
	loginButton = "/html/body/div[12]/div[1]/div[2]/div/div/div[1]/div[2]/div/div/div/form/div[2]/button"
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
		chromedp.SendKeys(usernameField, s.Credentials.Username),
		chromedp.SendKeys(passwordField, s.Credentials.Password),
		chromedp.Click(loginButton),
		//chromedp.Submit(passwordField),
	)

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible(logoutButton),
			Limit: 10 * time.Second, IsException: true, Message: "Login Failed"},
	)

	log.Info().Msgf("Successfully authenticated as: %v", s.Credentials.Username)
}
