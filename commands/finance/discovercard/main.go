package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/finance/discovercard"
	"github.com/walterjwhite/go-application/libraries/property"
)

var (
	s *discovercard.DiscoverCardSession
)

func init() {
	application.Configure()

	s = &discovercard.DiscoverCardSession{Credentials: &discovercard.WebCredentials{}}
	property.Load(s.Credentials, "")
}

func main() {
	s.GetBalance(application.Context)
}
