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

func (s *Session) authenticate(token string) error {
	token = s.trim(token)
	err := validateToken(token)
	if err != nil {
		return err
	}

	log.Info().Msg("session.authenticate - authenticating with provided credentials")

	ctx, cancel := context.WithTimeout(s.ctx, loginTimeout)
	defer cancel()

	err = action.Execute(ctx, chromedp.Navigate(s.Endpoint.Uri))
	if err != nil {
		return err
	}

	err = action.Execute(ctx,
		chromedp.SendKeys(USERNAME, strings.TrimSpace(s.Credentials.Domain+"\\"+s.Credentials.Username)),

		chromedp.SendKeys(PASSWORD, strings.TrimSpace(s.Credentials.Password)),
		chromedp.SendKeys(TOKEN, strings.TrimSpace(s.getTokenAndPin(token))),
	)
	if err != nil {
		return err
	}

	_, err = chromedp.RunResponse(ctx, chromedp.Click(LOGIN_BUTTON))
	if err != nil {
		return err
	}

	if !s.IsAuthenticated() {
		if s.GoogleProvider != nil {
			logging.Warn(s.GoogleProvider.PublishStatus("failed to authenticate", false), "authenticate.IsAuthenticated")
		}

		return errors.New("session.authenticate - failed to authenticate")
	}

	if s.GoogleProvider != nil {
		logging.Warn(s.GoogleProvider.PublishStatus("authenticated", true), "authenticate.authenticated")
	}

	return nil
}

func (s *Session) trim(token string) string {
	s.Credentials.Username = strings.TrimSpace(s.Credentials.Username)
	s.Credentials.Domain = strings.TrimSpace(s.Credentials.Domain)
	s.Credentials.Password = strings.TrimSpace(s.Credentials.Password)

	s.Credentials.Pin = strings.TrimSpace(s.Credentials.Pin)

	return strings.TrimSpace(token)
}

func validateToken(token string) error {
	if len(token) != 6 {
		return fmt.Errorf("invalid token format: expected 6-digit token")
	}

	for _, c := range token {
		if c < '0' || c > '9' {
			return fmt.Errorf("invalid token format: token must contain only digits")
		}
	}

	return nil
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
