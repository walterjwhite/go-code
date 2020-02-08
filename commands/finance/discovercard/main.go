package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/finance/discovercard"
	"github.com/walterjwhite/go-application/libraries/property"
)

var (
	s *discovercard.DiscoverSession
)

func init() {
	application.Configure()

	s = &discovercard.DiscoverSession{Credentials: &discovercard.WebCredentials{}}
	property.Load(s.Credentials, "")
}

func main() {
	s.Authenticate(application.Context)
}
