package gateway

import (
	"context"

	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/headless"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"
)

const (
	menuButtonXpath   = "//*[@id=\"userMenuBtn\"]/div"
	logoffButtonXpath = "//*[@id=\"menuLogOffBtn\"]"
)

type TokenProvider interface {
	ReadToken(ctx context.Context) *string
	OnSuccess(ctx context.Context)
	OnError(ctx context.Context, err error)
}

func (s *Session) InitializeChromeDP(ctx context.Context) {
	if !s.Headless {
		log.Warn().Msg("New remote session")
		s.session = remote.New(ctx)
	} else {
		log.Warn().Msg("New headless session")
		s.session = headless.New(ctx)
	}
}

func (s *Session) Authenticate(token string) {
	token = s.trim(token)
	validateToken(token)

	session.Execute(s.session, chromedp.Navigate(s.Endpoint.Uri))

	log.Debug().Msgf("username: %v", s.Credentials.Username)
	log.Debug().Msgf("domain: %v", s.Credentials.Domain)
	log.Debug().Msgf("password: %v", s.Credentials.Password)
	log.Debug().Msgf("pin/token: %v", s.getTokenAndPin(token))

	session.ExecuteWithDelay(s.session,
		delay.NewRandom(*s.Delay, *s.Delay),

		chromedp.SendKeys(s.Endpoint.UsernameXPath, strings.TrimSpace(s.Credentials.Domain+"\\"+s.Credentials.Username)),

		chromedp.SendKeys(s.Endpoint.PasswordXPath, strings.TrimSpace(s.Credentials.Password)),
		chromedp.SendKeys(s.Endpoint.TokenXPath, strings.TrimSpace(s.getTokenAndPin(token))),
	)

	_, err := chromedp.RunResponse(s.session.Context(), chromedp.Click(s.Endpoint.LoginButtonXPath))
	logging.Panic(err)
}

func (s *Session) getTokenAndPin(token string) string {
	return s.Credentials.Pin + token
}

func (s *Session) Logout() {
	session.Execute(s.session,
		chromedp.Click(menuButtonXpath),
		chromedp.Click(logoffButtonXpath),
	)
}

func (s *Session) IsAuthenticated() bool {
	if chromedpexecutor.Exists(s.session, time.Duration(time.Second*5), "userMenuBtn", chromedp.ByID) {
		log.Warn().Msg("user is authenticated - userMenuBtn is present")
		return true
	}

	citrixLightInstallButtonExists := chromedpexecutor.Exists(s.session, time.Duration(time.Second*5), "protocolhandler-welcome-installButton", chromedp.ByID)
	log.Warn().Msgf("user is authenticated - light install button: %v", citrixLightInstallButtonExists)

	return citrixLightInstallButtonExists
}
