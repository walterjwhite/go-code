package citrix

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

func (s *Session) PostLoad(ctx context.Context) {
	s.validate()

	for i := range s.Instances {
		logging.Error(s.Instances[i].PostLoad(ctx))
	}

	logging.Error(s.Worker.Validate())
}

func (s *Session) validate() {
	if len(s.Credentials.Domain) == 0 {
		logging.Error(errors.New("domain is required"))
	}
	if len(s.Credentials.Username) == 0 {
		logging.Error(errors.New("username is required"))
	}
	if len(s.Credentials.Password) == 0 {
		logging.Error(errors.New("password is required"))
	}
	if len(s.Credentials.Pin) == 0 {
		logging.Error(errors.New("pin is required"))
	}

	log.Info().Msg("session.Validate - session configuration is valid")
}

func (s *Session) Init(pctx context.Context) {
	log.Info().Msg("session.Init(pctx)")

	s.ctx, s.cancel = provider.New(s.Conf, pctx)

	s.Worker.Reset()
}
