package citrix

import (
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	_ "embed"
)

const (
	lockScreenMatchThreshold = 0.04
	minWidth                 = 100
	minHeight                = 100
)

var windows10StartButtonData []byte

func (i *Instance) isScreenLocked() bool {
	size := action.GetWindowSize(i.ctx)

	matches := i.match(lockScreenMatchThreshold, windows10StartButtonData, 0, float64(size.Height-minHeight), minWidth, minHeight)

	return len(matches) == 0
}
