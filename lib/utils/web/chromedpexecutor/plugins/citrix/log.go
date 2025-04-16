package citrix

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
)

func (i *Instance) captureScreenshot(name_template string) {
	if log.Debug().Enabled() {
		chromedpexecutor.FullScreenshot(i.ctx, fmt.Sprintf(name_template, i.Index))
	}
}
