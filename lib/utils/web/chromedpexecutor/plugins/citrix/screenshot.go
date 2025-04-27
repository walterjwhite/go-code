package citrix

import (
	"fmt"
	"github.com/andreyvit/locateimage"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"context"
)

func saveScreenshot(ctx context.Context, name_template string, args ...interface{}) {
	if log.Debug().Enabled() {
		action.FullScreenshot(ctx, fmt.Sprintf(name_template, args...))
	}
}

func (i *Instance) Match(matchThreshold float64, referenceImageData []byte, x, y, width, height float64) *locateimage.Match {
	return action.Match(i.ctx, matchThreshold, referenceImageData, x, y, width, height)
}
