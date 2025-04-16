package citrix

import (
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	_ "embed"
)

const (
	passwordPromptMatchThreshold   = 0.04
	passwordPromptWidthFromCenter  = 160
	passwordPromptHeightFromCenter = 300
)

var passwordPromptButtonData []byte

func (i *Instance) isPasswordPromptVisible() bool {
	size := action.GetWindowSize(i.ctx)

	matches := i.match(passwordPromptMatchThreshold, passwordPromptButtonData,
		float64(size.Width/2-passwordPromptWidthFromCenter), float64(size.Height/2-passwordPromptHeightFromCenter),
		passwordPromptWidthFromCenter*2, passwordPromptHeightFromCenter*2)

	return len(matches) == 1
}
