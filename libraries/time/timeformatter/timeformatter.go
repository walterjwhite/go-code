package timeformatter

import (
	"time"
)

type TimeFormatter interface {
	Format(time.Time) string
	Get() string
}

/*
func Get() string {
	return Format(time.Now())
}
*/
