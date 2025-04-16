package citrix

import (
	"github.com/rs/zerolog/log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

func (i *Instance) eliza() {
	log.Warn().Msgf("eliza is enabled - %d", i.Index)

	if !i.session.Headless {
		action.AttachMousePositionListener(i.ctx)
	}

	i.handleExternalMouseMovement()
	i.launchCommand("edge")


	i.moveMouse(100, 100)
	i.moveMouse(200, 100)
	i.moveMouse(200, 200)
	i.moveMouse(100, 200)

	i.captureScreenshot("/tmp/3.gateway-mouse-wiggle-%d.png")
}

func (i *Instance) launchCommand(commandName string) {
	logging.Panic(chromedp.Run(i.ctx, chromedp.KeyEvent(kb.Meta)))
	time.Sleep(time.Second)
	logging.Panic(chromedp.Run(i.ctx, chromedp.KeyEvent(commandName)))
	time.Sleep(time.Second)
	logging.Panic(chromedp.Run(i.ctx, chromedp.KeyEvent(kb.Enter)))
	time.Sleep(15 * time.Second)
}
