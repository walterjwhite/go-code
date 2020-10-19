package discovercard

import (

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor"
	"time"
)

const (
	logoutButton  = "/html/body/div/header/div[1]/span/a"
)

func (s *Session) Logout() {
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