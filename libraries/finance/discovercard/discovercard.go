package discovercard

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	url = "https://www.discover.com"
)

// TODO: generalize this ...
type WebCredentials struct {
	Username string
	Password string
}

type DiscoverCardSession struct {
	Credentials *WebCredentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}

func (c *WebCredentials) HasDefault() bool {
	return false
}

func (c *WebCredentials) Refreshable() bool {
	return false
}

func (c *WebCredentials) EncryptedFields() []string {
	return []string{"Username", "Password"}
}

func (s *DiscoverCardSession) Login(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.SendKeys("//*[@id=\"userid-content\"]", s.Credentials.Username),
		chromedp.SendKeys("//*[@id=\"password-content\"]", s.Credentials.Password),
	)
	
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.Click("//*[@id=\"log-in-button\"]"),
			Limit: 3*time.Second, IsException: false},
	)
	
	log.Info().Msg("Logged in")
	s.chromedpsession.Screenshot("logged.in.after.png")
}

func (s *DiscoverCardSession) GetBalance(ctx context.Context) {
	if s.chromedpsession == nil {
		s.Login(ctx)
	}

	defer s.Logout(ctx)
	s.navigateToCreditCardActivity(ctx)
}

func (s *DiscoverCardSession) navigateToCreditCardActivity(ctx context.Context) {
	s.chromedpsession.Screenshot("activity.before.png")
	
	log.Info().Msg("Fetching balance")
	
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible("#main-content", chromedp.ByID),
		Limit: 3*time.Second, IsException: false},
	)
	
	s.chromedpsession.Screenshot("activity.wait-visible.png")
	
	var innerHtml string
	innerHtml = ""
	
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.InnerHTML("//*[@id=\"main-content\"]/div[3]/div[1]/h2", &innerHtml),
		Limit: 3*time.Second, IsException: false},
	)
	
	log.Info().Msgf("Balance: %v", innerHtml)
	
	s.chromedpsession.Screenshot("activity.after.png")
}

func (s *DiscoverCardSession) Logout(ctx context.Context) {
	log.Info().Msg("Logging out")
	
	defer s.chromedpsession.Cancel()

	//body > div > header > div.navbar-header > span > a
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.Click("/html/body/div/header/div[1]/span/a"),
		Limit: 3*time.Second, IsException: false},
	)
	// depending on where we are within the site, the xpath also changes
	///html/body/div[1]/header/div/div/div[2]/div[2]/ul/li[6]/a
	
	log.Info().Msg("Logged out")
	s.chromedpsession.Screenshot("logged.out.after.png")
}
