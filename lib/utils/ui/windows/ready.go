package windows

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	pollInterval                  = 5 * time.Second
	termsAcceptanceRefreshTimeout = 2 * time.Second
	waitReadyTimeout              = 1 * time.Minute
)

func (c *WindowsConf) WaitReady(pctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(pctx, waitReadyTimeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()

		case <-ticker.C:
			log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible")
			visible, err := c.IsTermsAcceptanceButtonVisible(ctx)
			if err != nil {
				log.Warn().Msgf("windows.WaitReady.IsTermsAcceptanceButtonVisible - err - %v", err)
				return false, err
			}

			if visible {
				log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - visible")

				err = c.Controller.Click(ctx, 100, 100)
				err = c.Controller.Type(ctx, "\r")

				time.Sleep(termsAcceptanceRefreshTimeout)

				log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - click, type, sleep")
				continue
			}

			log.Debug().Msg("windows.WaitReady.IsTermsAcceptanceButtonVisible - not visible")
			visible, err = c.IsStartMenuButtonPresent(ctx)
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
