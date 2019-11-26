package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/finance/pnc"
	"github.com/walterjwhite/go-application/libraries/property"
)

var (
	s *pnc.PNCSession
)

func init() {
	application.Configure()

	s = &pnc.PNCSession{Credentials: &pnc.PNCCredentials{}}
	property.Load(s.Credentials, "")
}

func main() {
	s.GetBalance(application.Context)
}
