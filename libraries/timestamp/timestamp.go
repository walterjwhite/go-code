package timestamp

import (
	"fmt"
	"time"
)

const FORMAT = "%d.%d.%d.%d.%d.%d"

func Get() string {
	t := time.Now()

	return fmt.Sprintf(FORMAT, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
