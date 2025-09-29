package learning

import (
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

	"strings"
	"time"

	"context"

	"github.com/avast/retry-go"
	"github.com/walterjwhite/go-code/lib/utils/publisher"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/auth/microsoft"
)

const (
	signInButton = "/html/body/header/nav/div/a[1]"

	authInput  = "#auth-id-input"
	authButton = "#auth-id-button"

	linkedInLearningEnvironment = "/html/body/main/div/section/ul/li[1]/a/div"

	userMenuButton = "#hue-menu-trigger-ember21"
	logoutButton   = "//*[@id=\"ember23\"]/div/div"

	linkedInLearningLoginUrl = "https://www.linkedin.com/learning-login"

	existsTimeout = 1 * time.Second
)

func (s *Session) authenticate(publisher publisher.Publisher) {
	log.Info().Msgf("Session.authenticate - running with email: %v", s.EmailAddress)

	err := retry.Do(
		func() error {
			err := s.doTryAuthenticate(publisher)
			if err != nil {
				log.Warn().Err(err).Msg("Session.doTryAuthenticate - Error")
			}
			return err
		},
		retry.Attempts(s.AuthRetryAttempts),
		retry.Delay(s.AuthRetryDelay),
	)

	if err != nil {
		action.Screenshot(s.ctx, fmt.Sprintf("/tmp/linkedin-authenticate-error-%v.png", time.Now().Unix()))
		log.Fatal().Err(err).Msg("Session.authenticate - Error")
	}
}

func (s *Session) doTryAuthenticate(publisher publisher.Publisher) error {
	log.Debug().Msg("Session.doTryAuthenticate - start")

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	default:
	}

	err := s.fetchAuthenticationPage()
	if err != nil {
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - fetchAuthenticationPage - Error")
		return err
	}

	s.enterEmailAddress()
	if err != nil {
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - enterEmailAddress - Error")
		return err
	}

	s.selectEnvironment()
	if err != nil {
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - selectEnvironment - Error")
		return err
	}

	err = microsoft.Authenticate(s.ctx, s.EmailAddress, s.Password, publisher)
	if err != nil {
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - microsoftAuthenticate - Error")
		return err
	}

	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	err = chromedp.Run(ctx, chromedp.WaitReady(userMenuButton))
	if err != nil {
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - waitReady.userMenuButton - Error")
		return err
	}

	if !s.isAuthenticated() {
		err = errors.New("failed to authenticate")
		log.Warn().Err(err).Msg("Session.doTryAuthenticate - failed to authenticate")
		return err
	}

	log.Info().Msg("Session.doTryAuthenticate - successfully authenticated")
	return nil
}

func (s *Session) fetchAuthenticationPage() error {
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.Navigate(linkedInLearningLoginUrl),
		chromedp.WaitReady("body"))
}

func (s *Session) enterEmailAddress() error {
	log.Info().Msg("Session.enterEmailAddress - entering email-address")

	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.WaitReady(authInput),
		chromedp.SendKeys(authInput, strings.TrimSuffix(s.EmailAddress, "\n")),
		chromedp.Click(authButton))
}

func (s *Session) selectEnvironment() error {
	log.Info().Msg("Session.selectEnvironment - Clicking environment")
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.WaitReady(linkedInLearningEnvironment, chromedp.BySearch),
		chromedp.Click(linkedInLearningEnvironment, chromedp.BySearch))
}

func (s *Session) Logout() error {
	return action.Execute(s.ctx,
		chromedp.Click(userMenuButton),
		chromedp.Click(logoutButton))
}

func (s *Session) isAuthenticated() bool {
	if action.ExistsByCssSelector(s.ctx, userMenuButton) {
		log.Debug().Msg("Session.isAuthenticated - user is authenticated - userMenuBtn is present")
		return true
	}

	return false
}
