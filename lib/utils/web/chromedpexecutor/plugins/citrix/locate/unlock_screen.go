package citrix

import (
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	citrixWorkspaceMenuOffset                    = 110
	citrixWorkspaceMenuCtrlAltDeleteWidthOffset  = 50
	citrixWorkspaceMenuCtrlAltDeleteHeightOffset = 160
)

func (i *Instance) handleUnlockingScreen() {
	if i.isScreenLocked() {
		log.Warn().Msg("screen is locked")
		i.unlockScreen()
	} else {
		log.Info().Msg("screen is NOT locked")
	}
}

func (i *Instance) unlockScreen() {
	size := action.GetWindowSize(i.ctx)

	middleX := float64(size.Width / 2)
	timeBetweenSteps := 1 * time.Second

	logging.Panic(chromedp.Run(i.ctx,
		chromedp.MouseClickXY(middleX, 10),

		chromedp.Sleep(timeBetweenSteps),

		chromedp.MouseClickXY(middleX+citrixWorkspaceMenuOffset, 10),

		chromedp.Sleep(timeBetweenSteps),

		chromedp.MouseClickXY(middleX+citrixWorkspaceMenuOffset+citrixWorkspaceMenuCtrlAltDeleteWidthOffset, 10+citrixWorkspaceMenuCtrlAltDeleteHeightOffset),

		chromedp.Sleep(timeBetweenSteps),
	))

	if !i.isPasswordPromptVisible() {
		log.Warn().Msg("Password prompt is not visible, but we were expecting it to be visible")
		return
	}

	log.Info().Msg("Password prompt is visible, attempting to unlock")

	logging.Panic(chromedp.Run(i.ctx, chromedp.KeyEvent(i.session.Credentials.Password), chromedp.Sleep(timeBetweenSteps), chromedp.KeyEvent(kb.Enter)))

	log.Info().Msg("Sent password, should be unlocked")
}
