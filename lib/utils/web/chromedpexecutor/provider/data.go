package provider

import (
	"github.com/walterjwhite/go-code/lib/time/delay"
	"time"
)

type Conf struct {
	Headless         bool
	HeadlessViewport *HeadlessViewport
	Lightpanda       bool

	ProxyAddress string
	Remote       string

	Delay     time.Duration
	DelayType delay.DelayType
}

type HeadlessViewport struct {
	Width  int64
	Height int64
}
