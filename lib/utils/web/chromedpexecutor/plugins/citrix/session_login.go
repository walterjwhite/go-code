package citrix

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	USERNAME     = "#login"
	PASSWORD     = "#passwd"
	TOKEN        = "#passwd1"
	LOGIN_BUTTON = "#loginBtn"

	loginTimeout = 30 * time.Second

)

func (s *Session) authenticate(token string) {
	token = s.trim(token)
	validateToken(token)

	log.Info().Msgf("session.authenticate - authenticating with token: %v", token)
	log.Debug().Msgf("session.authenticate - credentials: %v | %v | %v | %v", s.Credentials.Username, s.Credentials.Domain, s.Credentials.Password, s.getTokenAndPin(token))

	ctx, cancel := context.WithTimeout(s.ctx, loginTimeout)
	defer cancel()

	logging.Panic(action.Execute(ctx, chromedp.Navigate(s.Endpoint.Uri)))

	logging.Panic(action.Execute(ctx,
		chromedp.SendKeys(USERNAME, strings.TrimSpace(s.Credentials.Domain+"\\"+s.Credentials.Username)),

		chromedp.SendKeys(PASSWORD, strings.TrimSpace(s.Credentials.Password)),
		chromedp.SendKeys(TOKEN, strings.TrimSpace(s.getTokenAndPin(token))),
	))

	_, err := chromedp.RunResponse(ctx, chromedp.Click(LOGIN_BUTTON))
	logging.Panic(err)

	if !s.IsAuthenticated() {
		logging.Warn(s.GoogleProvider.PublishStatus("failed to authenticate", false), "authenticate.IsAuthenticated")
		logging.Panic(errors.New("session.authenticate - failed to authenticate"))
	}

	logging.Warn(s.GoogleProvider.PublishStatus("authenticated", true), "authenticate.authenticated")
}

func (s *Session) trim(token string) string {
	s.Credentials.Username = strings.TrimSpace(s.Credentials.Username)
	s.Credentials.Domain = strings.TrimSpace(s.Credentials.Domain)
	s.Credentials.Password = strings.TrimSpace(s.Credentials.Password)

	s.Credentials.Pin = strings.TrimSpace(s.Credentials.Pin)

	return strings.TrimSpace(token)
}

func validateToken(token string) {
	if len(token) != 6 {
		logging.Panic(fmt.Errorf("please enter the 6-digit token: %v", token))
	}
}

func (s *Session) getTokenAndPin(token string) string {
	return s.Credentials.Pin + token
}

func (s *Session) IsAuthenticated() bool {
	if action.ExistsById(s.ctx, "userMenuBtn") {
		log.Debug().Msg("session.IsAuthenticated - user is authenticated - userMenuBtn is present")
		return true
	}

	citrixLightInstallButtonExists := action.ExistsById(s.ctx, "protocolhandler-welcome-installButton")
	log.Debug().Msgf("session.IsAuthenticated - user is authenticated - light install button: %v", citrixLightInstallButtonExists)

	return citrixLightInstallButtonExists
}
