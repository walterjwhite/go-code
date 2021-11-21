package writer

import (
	"github.com/walterjwhite/go-code/lib/external/spot/data"
)

type SpotWriter interface {
	Write(r *data.Record)
}
