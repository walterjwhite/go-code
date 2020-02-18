package discovercard

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	logoutButton  = "/html/body/div/header/div[1]/span/a"
	usernameField = "//*[@id=\"userid-content\"]"
	passwordField = "//*[@id=\"password-content\"]"
)

func (s *DiscoverSession) Authenticate(ctx context.Context) {
	if s.chromedpsession != nil {
		s.Logout()
	}

	s.chromedpsession = chromedpexecutor.LaunchRemoteBrowser(ctx)

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

func (s *DiscoverSession) Logout() {
	log.Info().Msg("Logging out")

	defer s.chromedpsession.Cancel()

	//body > div > header > div.navbar-header > span > a
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.Click(logoutButton),
			Limit: 3 * time.Second, IsException: true, Message: "Logout failed"},
	)
	// depending on where we are within the site, the xpath also changes
	///html/body/div[1]/header/div/div/div[2]/div[2]/ul/li[6]/a
}
