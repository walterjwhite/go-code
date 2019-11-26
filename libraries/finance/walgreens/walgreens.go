package discovercard

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	url = "https://walgreens.com"
)

// TODO: generalize this ...
type WalgreensCredentials struct {
	Username string
	Password string
}

type WalgreensSession struct {
	Credentials *WalgreensCredentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}

func (c *WalgreensCredentials) HasDefault() bool {
	return false
}

func (c *WalgreensCredentials) Refreshable() bool {
	return false
}

func (c *WalgreensCredentials) EncryptedFields() []string {
	return []string{"Username", "Password"}
}

func (s *WalgreensSession) Login(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.Click("//*[@id=\"signin-btn-header-2\"]"),
		chromedp.SendKeys("//*[@id=\"user-name\"]", s.Credentials.Username),
		chromedp.Click("//*[@id=\"submit_btn\"]"),
		chromedp.SendKeys("//*[@id=\"user_password\"]", s.Credentials.Password),
		chromedp.Click("//*[@id=\"submit_btn\"]"),
	)
	
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.Click("//*[@id=\"log-in-button\"]"),
			Limit: 3*time.Second, IsException: false},
	)
	
	log.Info().Msg("Logged in")
	s.chromedpsession.Screenshot("logged.in.after.png")
}

func (s *WalgreensSession) OrderPhotos(ctx context.Context) {
	if s.chromedpsession == nil {
		s.Login(ctx)
	}

	defer s.Logout(ctx)
	
	// //*[@id="photoOrg-addPhotos-qmp-btn"]
}

func (s *WalgreensSession) Logout(ctx context.Context) {
	log.Info().Msg("Logging out")
	
	defer s.chromedpsession.Cancel()

	//body > div > header > div.navbar-header > span > a
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.Click("//*[@id=\"wg-header-main\"]/div/div/div/div/div/div[2]/ul/li[1]/a"),
		Limit: 3*time.Second, IsException: false},
	)
	
	
	log.Info().Msg("Logged out")
	s.chromedpsession.Screenshot("logged.out.after.png")
}
