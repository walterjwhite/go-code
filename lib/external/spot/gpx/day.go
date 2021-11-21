package gpx

import (
	"time"

	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"github.com/walterjwhite/go-code/lib/time/timeformatter/day"

	"path"
)

func Day(s *data.Session, d *time.Time) []*data.Record {
	return get(path.Join(s.DataPath, day.Format(*d)))
}
