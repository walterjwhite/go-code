package vanguard

import (
	"github.com/walterjwhite/go-application/libraries/chromedpexecutor"
)

const (
	url = "https://www.vanguard.com"
)

// TODO: generalize this ...
type Credentials struct {
	Username string
	Password string
}

type VanguardSession struct {
	Credentials *Credentials

	chromedpsession *chromedpexecutor.ChromeDPSession
}
