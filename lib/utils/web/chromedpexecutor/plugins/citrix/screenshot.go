package citrix

import (
	"github.com/andreyvit/locateimage"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"image"
)

func (i *Instance) Match(matchThreshold float64, referenceImageData image.Image, x, y, width, height float64) *locateimage.Match {
	return action.Match(i.ctx, matchThreshold, referenceImageData, x, y, width, height)
}
