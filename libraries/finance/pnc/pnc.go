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
type Credentials struct {
	Username string
	Password string
}

type PNCSession struct {
	Credentials *Credentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}

func (s *PNCSession) GetBalance(ctx context.Context) {
	s.Authenticate(ctx)

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
