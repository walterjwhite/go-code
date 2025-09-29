package citrix

import (
	"context"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"time"
)

const (
	citrixActionButtonHalfHeight  = 20
	citrixActionButtonWidthOffset = 110
	citrixActionButtonCADHeightOffset = 186
	citrixActionButtonCADWidthOffset  = 50

	screenSizeTimeout         = 1 * time.Second
	citrixDelayBetweenActions = 500 * time.Millisecond

	ctrlAltDeleteTimeout     = 5 * time.Second
	postCtrlAltDeleteTimeout = 15 * time.Second
)

func (i *Instance) lock() error {
	if i.locked {
		log.Warn().Msgf("%v - Instance.lock - already locked", i)
		return nil
	}

	i.sendCtrlAltDelete()

	ctx, cancel := context.WithTimeout(i.ctx, postCtrlAltDeleteTimeout)
	defer cancel()

	time.Sleep(ctrlAltDeleteTimeout)

	err := chromedp.Run(ctx,
		chromedp.KeyEvent(kb.Enter))
	if err == nil {
		log.Info().Msgf("%v - Instance.lock - locked", i)
		i.locked = true
	}

	return err
}

func (i *Instance) unlock() error {
	if !i.locked {
		log.Warn().Msgf("%v - Instance.unlock - already unlocked", i)
		return nil
	}

	isLocked, err := i.WindowsConf.IsLocked(i.ctx)
	if err != nil {
		return err
	}

	if !isLocked {
		log.Warn().Msgf("%v - Instance.unlock - does not appear to be locked", i)
		i.locked = false
		return nil
	}

	log.Info().Msgf("%v - Instance.unlock - instance appears to be locked, unlocking", i)

	err = i.sendCtrlAltDelete()
	if err != nil {
		return err
	}

	time.Sleep(ctrlAltDeleteTimeout)

	ctx, cancel := context.WithTimeout(i.ctx, postCtrlAltDeleteTimeout)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.KeyEvent(i.session.Credentials.Password),
		chromedp.KeyEvent(kb.Enter))
	if err == nil {
		i.locked = false
	}

	return err
}

func (i *Instance) sendCtrlAltDelete() error {
	ctx, cancel := context.WithTimeout(i.ctx, ctrlAltDeleteTimeout)
	defer cancel()

	x, err := i.getCitrixActionButtonLocation()
	if err != nil {
		return err
	}

	log.Debug().Msgf("%v - Instance.sendCtrlAltDelete - located citrix action button: %f", i, x)

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
