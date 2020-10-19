package pnc

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor"
	"time"
)

func (s *Session) GetBalance(ctx context.Context) {
	//s.Authenticate(ctx)

	//defer s.Logout()

	s.navigateToCreditCardActivity(ctx)
}

func (s *Session) navigateToCreditCardActivity(ctx context.Context) {
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
