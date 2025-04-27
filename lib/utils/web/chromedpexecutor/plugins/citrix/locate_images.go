package citrix

import (
	_ "embed"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	matchThreshold = 0.04

	windows10UserProfileIconRadius = 96

	windows10AcceptTermsButtonHeight = 39
	windows10AcceptTermsButtonWidth  = 126

	passwordPromptWidthFromCenter  = 160
	passwordPromptHeightFromCenter = 300

	windows10StartButtonMinWidth  = 100
	windows10StartButtonMinHeight = 100
)

var (
	windows10LoggingInUserProfileIcon []byte

	windows10AcceptTermsButton []byte

	passwordPromptButtonData []byte

	windows10StartButtonData []byte
)

func (i *Instance) handlePrompts() {
	if i.isWindows10StartButtonVisible() {
		log.Warn().Msg("screen is not locked (windows icon is present)")
		return
	}

	screenCheckDuration := 5 * time.Second
	for iteration := 0; iteration < 10; iteration++ {
		if i.isWaitingForTermsAcceptance() {
			logging.Panic(chromedp.Run(i.ctx,
				chromedp.KeyEvent(kb.Enter)))
		} else if i.isLoggingIn() {
			log.Info().Msg("waiting for system to login")
		} else {
			log.Warn().Msg("neither waiting for terms acceptance or logging in")
			return
		}

		time.Sleep(screenCheckDuration)
	}

	log.Warn().Msg("Exited loop by exceeding count, not finding what we wanted")
}

func (i *Instance) isWindows10StartButtonVisible() bool {
	size := action.GetWindowSize(i.ctx)

	match := action.Match(i.ctx, matchThreshold, windows10StartButtonData, 0, float64(size.Height-windows10StartButtonMinHeight), windows10StartButtonMinWidth, windows10StartButtonMinHeight)

	return match != nil
}

func (i *Instance) isWaitingForTermsAcceptance() bool {
	size := action.GetWindowSize(i.ctx)

	match := action.Match(i.ctx, matchThreshold, windows10AcceptTermsButton,
		0, 0, float64(size.Width), float64(size.Height))

	return match != nil
}

func (i *Instance) isLoggingIn() bool {
	size := action.GetWindowSize(i.ctx)

	match := action.Match(i.ctx, matchThreshold, windows10LoggingInUserProfileIcon,
		float64(size.Width/2-windows10UserProfileIconRadius), float64(size.Height/2-windows10UserProfileIconRadius*4),
		windows10UserProfileIconRadius*2, windows10UserProfileIconRadius*8)

	return match != nil
}

func (i *Instance) isPasswordPromptVisible() bool {
	size := action.GetWindowSize(i.ctx)

	match := action.Match(i.ctx, matchThreshold, passwordPromptButtonData,
		float64(size.Width/2-passwordPromptWidthFromCenter), float64(size.Height/2-passwordPromptHeightFromCenter),
		passwordPromptWidthFromCenter*2, passwordPromptHeightFromCenter*2)

	return match != nil
}
