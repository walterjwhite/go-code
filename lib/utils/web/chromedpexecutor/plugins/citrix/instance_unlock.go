package citrix

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/rs/zerolog/log"

	"time"
)

func (i *Instance) unlock() error {
	if !i.locked {
		log.Debug().Msgf("%v - Instance.unlock - already unlocked", i)
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
