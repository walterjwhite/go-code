package citrix

import (
	"context"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

func (i *Instance) sendCtrlAltDelete() error {
	x, err := i.getCitrixActionButtonLocation()
	if err != nil {
		log.Warn().Msgf("%v - Instance.sendCtrlAltDelete - error locating citrix action button", i)
		return err
	}

	log.Info().Msgf("%v - Instance.sendCtrlAltDelete - located citrix action button: %f", i, x)

	ctx, cancel := context.WithTimeout(i.ctx, ctrlAltDeleteTimeout)
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.MouseEvent(input.MouseMoved, x, citrixActionButtonHalfHeight),
		chromedp.Sleep(citrixDelayBetweenActions),

		chromedp.MouseClickXY(x, citrixActionButtonHalfHeight),
		chromedp.Sleep(citrixDelayBetweenActions),

		chromedp.MouseClickXY(x+citrixActionButtonWidthOffset, citrixActionButtonHalfHeight),
		chromedp.Sleep(citrixDelayBetweenActions),

		chromedp.MouseClickXY(x+citrixActionButtonWidthOffset+citrixActionButtonCADWidthOffset, citrixActionButtonHalfHeight+citrixActionButtonCADHeightOffset),
	)
}

func (i *Instance) getCitrixActionButtonLocation() (float64, error) {
	ctx, cancel := context.WithTimeout(i.ctx, screenSizeTimeout)
	defer cancel()

	size, err := action.GetWindowSize(ctx)
	if err != nil {
		return 0, err
	}

	return float64(size.Width) / 2.0, nil
}
