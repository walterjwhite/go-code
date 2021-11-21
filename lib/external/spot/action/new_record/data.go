package new_record

import (
	"context"

	"github.com/walterjwhite/go-code/lib/external/spot/data"
)

type Configuration struct {
	Session *data.Session
}

func New(s *data.Session) *Configuration {
	return &Configuration{Session: s}
}

func (c *Configuration) Init(s *data.Session, ctx context.Context) {
}
