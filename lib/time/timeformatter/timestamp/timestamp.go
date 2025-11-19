package timestamp

import (
	"fmt"

	"time"
)

type Configuration struct {
	Template string
}

var (
	Default *Configuration
)

func init() {
	Default = &Configuration{Template: "%d.%d.%d.%d.%d.%d.%d"}
}

func UseNested() {
	Default = &Configuration{Template: "%d/%02d.%s/%d/%d.%d.%d"}
}

func (c *Configuration) Format(t time.Time) string {
	if c.Template == "%d/%02d.%s/%d/%d.%d.%d" {
		return fmt.Sprintf(c.Template, t.Year(), t.Month(), t.Month().String(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
	return fmt.Sprintf(c.Template, t.Year(), t.Month(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
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
