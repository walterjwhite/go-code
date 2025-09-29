package microsoft

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/utils/publisher"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	microsoftEmailInputId                 = "#i0116"
	microsoftAuthenticationSubmitButtonId = "#idSIButton9"

	microsoftPasswordInputId = "#i0118"

	microsoft2FAId            = "#idRichContext_DisplaySign"
	microsoft2FARefreshPeriod = "//*[@id=\"idLbl_SAOTCAS_TD_Cb\"]/span"

	stepTimeout      = 15 * time.Second
	extractTimeout   = 250 * time.Millisecond
	twoFactorTimeout = 1 * time.Minute

	microsoftDelay = 1 * time.Second
)

func Authenticate(pctx context.Context, emailAddress, password string, publisher publisher.Publisher) error {
	log.Info().Msg("performing microsoft authentication")

	log.Debug().Msgf("using fixed delay [for microsoft]: %v", microsoftDelay)
	mpctx := context.WithValue(pctx, action.ContextKey, delay.New(microsoftDelay))

	err := enterEmail(mpctx, emailAddress)
	if err != nil {
		return err
	}

	err = enterPassword(mpctx, password)
	if err != nil {
		return err
	}

	err = twoFactor(mpctx, publisher)
	if err != nil {
		return err
	}

	log.Info().Msg("completed microsoft authentication")
	return nil
}

func enterEmail(pctx context.Context, emailAddress string) error {
	log.Info().Msg("entering email address")
	ctx, cancel := context.WithTimeout(pctx, stepTimeout)
	defer cancel()

	err := action.Execute(ctx,
		chromedp.WaitReady(microsoftEmailInputId),
		chromedp.WaitReady(microsoftAuthenticationSubmitButtonId),

		chromedp.SendKeys(microsoftEmailInputId, strings.TrimSuffix(emailAddress, "\n")),

		chromedp.Click(microsoftAuthenticationSubmitButtonId))
	if err != nil {
		return err
	}

	log.Info().Msg("entered email address")
	return nil
}

func enterPassword(pctx context.Context, password string) error {
	ctx, cancel := context.WithTimeout(pctx, stepTimeout)
	defer cancel()

	log.Info().Msg("waiting for password input to be ready")
	err := action.Execute(ctx,
		chromedp.WaitReady(microsoftPasswordInputId),
		chromedp.WaitReady(microsoftAuthenticationSubmitButtonId),
		chromedp.SendKeys(microsoftPasswordInputId, strings.TrimSuffix(password, "\n")),
		chromedp.Click(microsoftAuthenticationSubmitButtonId))
	if err != nil {
		return err
	}

	log.Info().Msg("entered password")
	return nil
}

func twoFactor(pctx context.Context, publisher publisher.Publisher) error {
	log.Info().Msg("getting token")

	ctx, cancel := context.WithTimeout(pctx, stepTimeout)
	defer cancel()

	action.Screenshot(pctx, fmt.Sprintf("/tmp/microsoft-auth-before-wait-ready-%v.png", time.Now().Unix()))
	err := action.Execute(ctx,
		chromedp.WaitReady(microsoft2FAId))
	if err != nil {
		action.Screenshot(pctx, fmt.Sprintf("/tmp/microsoft-auth-error-wait-ready-%v.png", time.Now().Unix()))
		return err
	}

	token, err := action.Get(ctx, microsoft2FAId)
	if err != nil {
		action.Screenshot(pctx, fmt.Sprintf("/tmp/microsoft-auth-error-%v.png", time.Now().Unix()))
		return err
	}

	log.Info().Msgf("Microsoft token: %s", token)
	if publisher != nil {
		logging.Warn(publisher.Publish([]byte(fmt.Sprintf("microsoft token: %s", token))), false, "microsoft.twoFactor.publisher.Publish")
	}

	ctx2, cancel2 := context.WithTimeout(pctx, stepTimeout)
	defer cancel2()

	err = action.Execute(ctx2,
		chromedp.WaitReady(microsoft2FARefreshPeriod),
		chromedp.Click(microsoft2FARefreshPeriod),
		chromedp.WaitReady("body"))
	if err != nil {
		return err
	}

	timeout := time.After(twoFactorTimeout)
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return errors.New("timed out after 1m - 2FA token is required")
		case <-ticker.C:
			if !action.ExistsByCssSelector(pctx, microsoft2FAId) {
				return nil
			}
		}
	}
}
