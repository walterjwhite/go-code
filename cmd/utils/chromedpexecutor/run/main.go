package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	log.Info().Msg("before")

	run.Run(application.Context)

	log.Info().Msg("after")
}
