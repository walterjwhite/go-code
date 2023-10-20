package run

import (
	"context"
	"path/filepath"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/property"
)

func New(path string, applications ...string) *Session {
	s := &Session{Path: path}

	s.Applications = make([]*Application, len(applications))
	for index, application := range applications {
		a := &Application{}

		appConfPath := filepath.Join(s.Path, application, ".application", "go.yaml")
		property.LoadFileWithPath(a, appConfPath)
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
	logging.Panic(err)
}
