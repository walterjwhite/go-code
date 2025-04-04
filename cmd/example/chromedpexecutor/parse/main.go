package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
)

func main() {
	log.Info().Msgf("parsed: %v", run.ParseAction("key,\"\u010e\""))
}
