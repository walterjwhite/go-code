package windows

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	pollInterval                  = 5 * time.Second
	termsAcceptanceRefreshTimeout = 2 * time.Second
	waitReadyTimeout              = 1 * time.Minute
	maxRetryCount                 = 12 // Maximum number of polling attempts
)

func (c *WindowsConf) WaitReady(pctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(pctx, waitReadyTimeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()

		case <-ticker.C:
			retryCount++
			if retryCount > maxRetryCount {
				return false, errors.New("maximum retry count exceeded while waiting for windows ready state")
			}

			log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible")
			visible, err := c.IsTermsAcceptanceButtonVisible(ctx)
			if err != nil {
				log.Warn().Msgf("windows.WaitReady.IsTermsAcceptanceButtonVisible - err - %v", err)
				return false, err
			}

			if visible {
				log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - visible")

				clickCtx, clickCancel := context.WithTimeout(ctx, 10*time.Second)
				err = c.Controller.Click(clickCtx, 100, 100)
				clickCancel()
				if err != nil {
					log.Warn().Msgf("windows.WaitReady.Click - err - %v", err)
					return false, err
				}

				typeCtx, typeCancel := context.WithTimeout(ctx, 10*time.Second)
				err = c.Controller.Type(typeCtx, "\r")
				typeCancel()
				if err != nil {
					log.Warn().Msgf("windows.WaitReady.Type - err - %v", err)
					return false, err
				}
				time.Sleep(termsAcceptanceRefreshTimeout)

				log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - click, type, sleep")
				continue
			}

			log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - not visible")

			checkCtx, checkCancel := context.WithTimeout(ctx, 10*time.Second)
			visible, err = c.IsStartMenuButtonPresent(checkCtx)
			checkCancel()
			if err != nil {
				log.Warn().Msgf("windows.WaitReady.StartMenuButtonPresent - err - %v", err)
				return false, err
			}

			if visible {
				log.Debug().Msg("windows.WaitReady.StartMenuButtonPresent - visible")
				return true, nil
			}

			log.Debug().Msg("windows.WaitReady.StartMenuButtonPresent - not visible")
		}
	}
}
