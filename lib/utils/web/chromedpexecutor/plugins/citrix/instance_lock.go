package citrix

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/rs/zerolog/log"

	"time"
)

const (
	citrixActionButtonHalfHeight  = 20
	citrixActionButtonWidthOffset = 110
	citrixActionButtonCADHeightOffset = 160
	citrixActionButtonCADWidthOffset = 50

	screenSizeTimeout         = 1 * time.Second
	citrixDelayBetweenActions = 500 * time.Millisecond

	ctrlAltDeleteTimeout     = 5 * time.Second
	postCtrlAltDeleteTimeout = 15 * time.Second
)

func (i *Instance) lock() error {
	if !i.Lockable {
		log.Warn().Msgf("%v - Instance.lock - instance cannot be locked", i)
		return nil
	}

	if i.locked {
		log.Warn().Msgf("%v - Instance.lock - already locked", i)
		return nil
	}

	err := i.sendCtrlAltDelete()
	if err != nil {
		log.Info().Msgf("%v - error sending ctrl+alt+delete", i)
		return err
	}

	time.Sleep(ctrlAltDeleteTimeout)

	ctx, cancel := context.WithTimeout(i.ctx, postCtrlAltDeleteTimeout)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.KeyEvent(kb.Enter))
	if err == nil {
		log.Info().Msgf("%v - Instance.lock - locked", i)
		i.locked = true
	}

	return err
}
