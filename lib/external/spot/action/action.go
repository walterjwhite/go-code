package action

import (
	"context"

	"github.com/walterjwhite/go-code/lib/external/spot/data"
)

type RecordAction interface {
	OnNewRecord(old, new *data.Record)
	//Init(latestRecord *data.Record, ctx context.Context)
}

type BackgroundAction interface {
	Init(s *data.Session, ctx context.Context)
}
