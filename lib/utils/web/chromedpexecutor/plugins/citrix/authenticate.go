package citrix

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	usernameXPath    = "//*[@id=\"login\"]"
	passwordXPath    = "//*[@id=\"passwd\"]"
	tokenXPath       = "//*[@id=\"passwd1\"]"
	loginButtonXPath = "//*[@id=\"loginBtn\"]"

	menuButtonXpath   = "//*[@id=\"userMenuBtn\"]/div"
	logoffButtonXpath = "//*[@id=\"menuLogOffBtn\"]"

	loginTimeout  = 30 * time.Second
	logoutTimeout = 5 * time.Second
	existsTimeout = 1 * time.Second
)

func (s *Session) authenticate(token string) {
	token = s.trim(token)
	validateToken(token)

	log.Info().Msgf("Session.authenticate - authenticating with token: %v", token)
	log.Debug().Msgf("Session.authenticate - credentials: %v | %v | %v | %v", s.Credentials.Username, s.Credentials.Domain, s.Credentials.Password, s.getTokenAndPin(token))

	ctx, cancel := context.WithTimeout(s.ctx, loginTimeout)
	defer cancel()

	action.Execute(ctx, chromedp.Navigate(s.Endpoint.Uri))

	action.Execute(ctx,
		chromedp.SendKeys(usernameXPath, strings.TrimSpace(s.Credentials.Domain+"\\"+s.Credentials.Username)),

		chromedp.SendKeys(passwordXPath, strings.TrimSpace(s.Credentials.Password)),
		chromedp.SendKeys(tokenXPath, strings.TrimSpace(s.getTokenAndPin(token))),
	)

	_, err := chromedp.RunResponse(ctx, chromedp.Click(loginButtonXPath))
	logging.Panic(err)

	if !s.IsAuthenticated() {
		s.GoogleProvider.PublishStatus("failed to authenticate", false)
		logging.Panic(errors.New("Session.authenticate - failed to authenticate"))
	}

	s.GoogleProvider.PublishStatus("authenticated", true)
}

func (s *Session) getTokenAndPin(token string) string {
	return s.Credentials.Pin + token
}

func (s *Session) Logout() {
	ctx, cancel := context.WithTimeout(s.ctx, logoutTimeout)
	defer cancel()

	action.Execute(ctx,
		chromedp.Click(menuButtonXpath),
		chromedp.Click(logoffButtonXpath),
	)
}

func (s *Session) IsAuthenticated() bool {
	if action.ExistsById(s.ctx, "userMenuBtn") {
		log.Debug().Msg("Session.IsAuthenticated - user is authenticated - userMenuBtn is present")
		return true
	}

	citrixLightInstallButtonExists := action.ExistsById(s.ctx, "protocolhandler-welcome-installButton")
	log.Debug().Msgf("Session.IsAuthenticated - user is authenticated - light install button: %v", citrixLightInstallButtonExists)

	return citrixLightInstallButtonExists
}
