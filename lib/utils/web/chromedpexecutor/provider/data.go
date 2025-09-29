package provider

import (
	"github.com/walterjwhite/go-code/lib/time/delay"
	"time"
)

type Conf struct {
	Headless     bool
	ProxyAddress string
	Remote       string

	Delay     time.Duration
	DelayType delay.DelayType
}
