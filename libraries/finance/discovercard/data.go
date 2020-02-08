package discovercard

import (
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
)

// TODO: generalize this ...
type WebCredentials struct {
	Username string
	Password string
}

// TODO: add support for entering one-time tokens
// TODO: add support for answering challenge questions

type DiscoverSession struct {
	Credentials *WebCredentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}
