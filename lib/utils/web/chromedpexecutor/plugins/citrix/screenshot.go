package citrix

import (
	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"image"
)

func (i *Instance) Match(matchThreshold float64, referenceImageData image.Image, x, y, width, height float64) (bool, error) {
	m := graphical.ImageMatch{Ctx: i.ctx, Image: referenceImageData,
		MatchRegion:    &graphical.MatchRegion{X: x, Y: y, Width: width, Height: height},
		MatchThreshold: matchThreshold,
		Controller:     i.WindowsConf.Controller}

	return m.Matches()
}
