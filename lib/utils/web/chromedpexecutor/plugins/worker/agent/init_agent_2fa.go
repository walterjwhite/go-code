package agent

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/until"
	"github.com/walterjwhite/go-code/lib/utils/ocr"
	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"

	"time"
)

func (c *Conf) waitFor2FAToComplete(pctx context.Context) error {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	log.Info().Msg("agent.waitFor2FAToComplete.2FA - checking if 2FA is present")

	i := graphical.ImageMatch{Ctx: ctx, Image: microsoft2FAHeaderImage, Controller: c.WindowsConf.Controller}
	err := until.Until(ctx, 250*time.Millisecond, i.Matches)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Warn().Msg("timed out waiting for 2FA, must not be present on screen")
			return nil
		}

		log.Warn().Msgf("agent.waitFor2FAToComplete.is2FAPresent.error - %v", err)
		return err
	}

	log.Info().Msg("agent.waitFor2FAToComplete.2FA - is present")

	tokenImageData, err := c.WindowsConf.Controller.ScreenshotOf(ctx, float64(i.Match.Rect.Min.X), float64(i.Match.Rect.Max.Y), float64(i.Match.Rect.Max.X-i.Match.Rect.Min.X), 100)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, err = graphical.ImageToBytes(tokenImageData)
	if err != nil {
		return err
	}

	token, err := ocr.Text(bytes)
	if err != nil {
		logging.Warn(err, false, "waitFor2FAToComplete.ocr.Text")
		return err
	}

	log.Info().Msgf("token: %s", token)





	ctx, cancel = context.WithTimeout(pctx, 1*time.Minute)
	defer cancel()

	err = until.Until(ctx, 250*time.Millisecond, i.NotMatches)
	if err != nil {
		log.Warn().Msgf("waitFor2FAToComplete.is2FANotPresent.error - %v", err)
		return err
	}

	log.Info().Msg("waitFor2FAToComplete.2FA - is no longer present")
	return nil
}
