package decade

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"strconv"
	"time"
)

type Configuration struct {
	Template string
}

var (
	Default *Configuration
)

func init() {
	Default = &Configuration{Template: "%v-%v"}
}

func (c *Configuration) Format(t time.Time) string {
	y := t.Year()
	year := strconv.Itoa(y)

	last := year[len(year)-1:]
	allButLast := year[:len(year)-1]

	var start, end int
	if last == "1" {
		start, end = atStart(y)
	} else if last == "0" {
		start, end = atEnd(y)
	} else {
		nstart, err := strconv.Atoi(allButLast + "1")
		logging.Panic(err)

		start, end = atStart(nstart)
	}

	return fmt.Sprintf(c.Template, start, end)
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

func atStart(y int) (int, int) {
	return y, y + 9
}

func atEnd(y int) (int, int) {
	return y - 9, y
}
