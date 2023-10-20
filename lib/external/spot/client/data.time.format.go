package client

import (
	"fmt"
	"strings"
	"time"
)

type SpotTime time.Time

const ctLayout = "2006-01-02T00:00:00-0700"

func (ct *SpotTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = SpotTime(nt)
	return
}

func (ct SpotTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *SpotTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}

func (ct *SpotTime) Time() time.Time {
	return time.Time(*ct)
}
