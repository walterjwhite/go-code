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

func (s *DiscoverSession) GetBalance(ctx context.Context) {
	s.Authenticate(ctx)

	defer s.Logout(ctx)
	s.navigateToCreditCardActivity(ctx)
}

func (s *DiscoverSession) navigateToCreditCardActivity(ctx context.Context) {
	s.chromedpsession.Screenshot("activity.before.png")

	log.Info().Msg("Fetching balance")

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.WaitVisible("#main-content", chromedp.ByID),
			Limit: 3 * time.Second, IsException: false},
	)

	s.chromedpsession.Screenshot("activity.wait-visible.png")

	innerHtml := ""

	s.chromedpsession.ExecuteTimeLimited(
		chromedpexecutor.TimeLimitedChromeAction{Action: chromedp.InnerHTML("//*[@id=\"main-content\"]/div[3]/div[1]/h2", &innerHtml),
			Limit: 3 * time.Second, IsException: false},
	)

	log.Info().Msgf("Balance: %v", innerHtml)

	s.chromedpsession.Screenshot("activity.after.png")
}
