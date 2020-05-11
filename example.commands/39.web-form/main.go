package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/token/plugins/web"
)

func init() {
	application.Configure()
}

// TODO:
// record approval / denial (to file, to ES)
// dynamically serve requests (set timeout for approval, assume denied if not approved by specified time, allow approval to be denied within a given time frame)
// record (client IP address, browser, other headers)
// DONE
// 1. inject request #
// 2. inject request description
func main() {
	webReader := web.NewRandomWebReader()
	webReader.Context = application.Context
	webReader.Cancel = application.Cancel

	log.Info().Msgf("GOT: %v", webReader.Get())
}
