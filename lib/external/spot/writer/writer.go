package writer

import (
	"github.com/walterjwhite/go/lib/external/spot/data"
)

type SpotWriter interface {
	Write(r *data.Record)
}
