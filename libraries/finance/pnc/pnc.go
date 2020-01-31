package pnc

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
	"time"
)

const (
	url = "https://www.pnc.com"
)

// TODO: generalize this ...
type PNCCredentials struct {
	Username string
	Password string
}

type PNCSession struct {
	Credentials *PNCCredentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}

func (c *PNCCredentials) EncryptedFields() []string {
	return []string{"Username", "Password"}
}

func (s *PNCSession) Login(ctx context.Context) {
	s.chromedpsession = chromedpexecutor.New(ctx)

	//defer s.Cancel()

	s.chromedpsession.Execute(
		chromedp.Navigate(url),
		chromedp.SendKeys("//*[@id=\"userId\"]", s.Credentials.Username),
		chromedp.SendKeys("//*[@id=\"passwordInputField\"]", s.Credentials.Password),
		chromedp.Click("//*[@id=\"olb-btn\"]"),
	)
}

func (s *PNCSession) Logout() {
	defer s.chromedpsession.Cancel()

	s.chromedpsession.Execute(chromedp.Click("//*[@id=\"topLinks\"]/ul/li[3]/a"))
}

func (s *PNCSession) GetBalance(ctx context.Context) {
	if s.chromedpsession == nil {
		s.Login(ctx)
	}

	defer s.Logout()

	s.navigateToCreditCardActivity(ctx)
}

func (s *PNCSession) navigateToCreditCardActivity(ctx context.Context) {
	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible("//*[@id=\"cc_detail_row\"]/td[4]"),
			Limit: 3 * time.Second, IsException: false},
	)

	innerHtml := ""

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.InnerHTML("//*[@id=\"cc_detail_row\"]/td[4]", &innerHtml),
			Limit: 3 * time.Second, IsException: false},
	)

	log.Info().Msgf("Balance: %v", innerHtml)
}
