package daily_export

import (
	"context"

	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"github.com/walterjwhite/go-code/lib/time/periodic"
	"time"
)

type DailyExportConfiguration struct {
	Session *data.Session

	periodicInstance *periodic.PeriodicInstance
}

func New(s *data.Session) *DailyExportConfiguration {
	return &DailyExportConfiguration{Session: s}
}

func (c *DailyExportConfiguration) Init(s *data.Session, ctx context.Context) {
	t := 24 * time.Hour
	c.periodicInstance = periodic.Periodic(ctx, &t, false, c.doExport)
}
