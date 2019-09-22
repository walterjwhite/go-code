package server

import (
	"github.com/walterjwhite/go-application/libraries/os/unix/tail"
)

var logData string

func (s *Server) Logs(args *Args, response *string) error {
	*response = logData
	return nil
}

func RefreshLogs() {
	logData = tail.Data()
}
