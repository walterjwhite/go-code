package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/finance/vanguard"
	"github.com/walterjwhite/go-application/libraries/property"
)

var (
	s *vanguard.VanguardSession
)

func init() {
	application.Configure()

	s = &vanguard.VanguardSession{Credentials: &vanguard.Credentials{}}
	property.Load(s.Credentials, "")
}

func main() {
	s.Authenticate(application.Context)
}
