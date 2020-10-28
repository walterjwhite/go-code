package discovercard

import (
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
)

// TODO: generalize this ...
type Credentials struct {
	Username string
	Password string
}

// TODO: add support for entering one-time tokens
// TODO: add support for answering challenge questions

type Session struct {
	Credentials *Credentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}
