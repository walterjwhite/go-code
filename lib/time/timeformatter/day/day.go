package day

import (

	"github.com/walterjwhite/go-code/lib/time/timeformatter/decade"

	"time"
)

type Configuration struct {
	Template string
}

var (
	Default *Configuration
)

func init() {
	Default = &Configuration{Template: "2006/01.January/02"}
}

func (c *Configuration) Format(t time.Time) string {
	d := decade.Format(t)

	return d + "/" + t.Format(c.Template)
}

func (c *Configuration) Get() string {
	return c.Format(time.Now())
}

func Format(t time.Time) string {
	return Default.Format(t)
}

func Get() string {
	return Format(time.Now())
}
