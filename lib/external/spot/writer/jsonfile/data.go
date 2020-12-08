package jsonfile

import (
	"github.com/walterjwhite/go/lib/external/spot/data"
)

type RecordAppenderConfiguration struct {
	Session *data.Session
}

func New(s *data.Session) *RecordAppenderConfiguration {
	return &RecordAppenderConfiguration{Session: s}
}
