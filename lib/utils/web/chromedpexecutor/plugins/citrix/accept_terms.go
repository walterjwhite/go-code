package citrix

import (
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	_ "embed"
)

const (
	windows10AcceptTermsThreshold    = 0.04
	windows10AcceptTermsButtonHeight = 39
	windows10AcceptTermsButtonWidth  = 126
)

var windows10AcceptTermsButton []byte

func (i *Instance) isWaitingForTermsAcceptance() bool {
	size := action.GetWindowSize(i.ctx)

	matches := i.match(lockScreenMatchThreshold, windows10AcceptTermsButton,
		float64(size.Width/2-windows10AcceptTermsButtonWidth/2), float64(size.Height/2-windows10AcceptTermsButtonHeight*4),
		windows10AcceptTermsButtonWidth, windows10AcceptTermsButtonHeight*4)

	return len(matches) == 1
}
