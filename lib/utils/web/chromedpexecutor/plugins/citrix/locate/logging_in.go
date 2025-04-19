package citrix

import (
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	_ "embed"
)

const (
	windows10LoggingInUserProfileMatchThreshold = 0.04
	windows10UserProfileIconRadius              = 96
)

var windows10LoggingInUserProfileIcon []byte

func (i *Instance) isLoggingIn() bool {
	size := action.GetWindowSize(i.ctx)

	matches := i.match(windows10LoggingInUserProfileMatchThreshold, windows10LoggingInUserProfileIcon,
		float64(size.Width/2-windows10UserProfileIconRadius), float64(size.Height/2-windows10UserProfileIconRadius*4),
		windows10UserProfileIconRadius*2, windows10UserProfileIconRadius*8)

	return len(matches) == 1
}
