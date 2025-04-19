package citrix

import (
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	menuButtonXpath   = "//*[@id=\"userMenuBtn\"]/div"
	logoffButtonXpath = "//*[@id=\"menuLogOffBtn\"]"
)

func (s *Session) authenticate(token string) {
	token = s.trim(token)
	validateToken(token)

	log.Info().Msgf("running with: %v", token)

	action.Execute(s.ctx, chromedp.Navigate(s.Endpoint.Uri))

	log.Debug().Msgf("username: %v", s.Credentials.Username)
	log.Debug().Msgf("domain: %v", s.Credentials.Domain)
	log.Debug().Msgf("password: %v", s.Credentials.Password)
	log.Debug().Msgf("pin/token: %v", s.getTokenAndPin(token))

	action.ExecuteWithDelay(s.ctx,
		delay.NewRandom(*s.Delay, *s.Delay),

		chromedp.SendKeys(s.Endpoint.UsernameXPath, strings.TrimSpace(s.Credentials.Domain+"\\"+s.Credentials.Username)),

		chromedp.SendKeys(s.Endpoint.PasswordXPath, strings.TrimSpace(s.Credentials.Password)),
		chromedp.SendKeys(s.Endpoint.TokenXPath, strings.TrimSpace(s.getTokenAndPin(token))),
	)

	_, err := chromedp.RunResponse(s.ctx, chromedp.Click(s.Endpoint.LoginButtonXPath))
	logging.Panic(err)
}

func (s *Session) getTokenAndPin(token string) string {
	return s.Credentials.Pin + token
}

func (s *Session) logout() {
	action.Execute(s.ctx,
		chromedp.Click(menuButtonXpath),
		chromedp.Click(logoffButtonXpath),
	)
}

func (s *Session) isAuthenticated() bool {
	if action.Exists(s.ctx, time.Duration(time.Second*5), "userMenuBtn", chromedp.ByID) {
		log.Warn().Msg("user is authenticated - userMenuBtn is present")
		return true
	}

	citrixLightInstallButtonExists := action.Exists(s.ctx, time.Duration(time.Second*5), "protocolhandler-welcome-installButton", chromedp.ByID)
	log.Warn().Msgf("user is authenticated - light install button: %v", citrixLightInstallButtonExists)

	return citrixLightInstallButtonExists
}
