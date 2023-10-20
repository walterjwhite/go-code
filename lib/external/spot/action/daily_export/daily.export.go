package daily_export

import (
	"github.com/walterjwhite/go-code/lib/external/spot/gpx"
	"github.com/walterjwhite/go-code/lib/time/timeformatter/day"

	"path/filepath"
)

func (c *DailyExportConfiguration) doExport() error {
	gpx.Export(gpx.Latest(c.Session), filepath.Join(c.Session.SessionPath, ".exports", day.Get()+".gpx"))
	return nil
}
