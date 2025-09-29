package citrix

import (
	"github.com/rs/zerolog/log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

func (i *Instance) Work() {
	log.Warn().Msgf("Instance.Work - eliza is enabled - %d", i.Index)

}
