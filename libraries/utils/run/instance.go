package run

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/application/property"
	"path/filepath"
)

func New(path string, applications ...string) *Session {
	s := &Session{Path: path}

	s.Applications = make([]*Application, len(applications))
	for index, application := range applications {
		a := &Application{}

		appConf := &property.Configuration{Path: filepath.Join(s.Path, application, "application.yaml")}
		appConf.Load(a)
		a.session = s

		s.Applications[index] = a
	}

	return s
}

func (s *Session) Run(ctx context.Context) {
	for _, a := range s.Applications {
		a.Run(ctx)
	}

	s.waitForAll()
}

func (s *Session) waitForAll() {
	for _, a := range s.Applications {
		_, err := a.command.Process.Wait()

		s.onError(err, a)
	}
}

func (s *Session) onError(err error, a *Application) {
	// TODO: push event to channel
	logging.Panic(err)
}
