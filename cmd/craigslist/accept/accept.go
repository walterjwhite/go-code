package main

import (
	"context"
	"flag"

	"time"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/foreachfile"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/craigslist"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

var (
	craigslistPostLinkDirectoryRootFlag = flag.String("CraigslistPostLinkDirectoryRoot", "./", "Root directory containing craigslist post links")
	craigslistPostLinkAcceptTimeoutFlag = flag.String("CraigslistPostLinkAcceptTimeout", "5s", "Timeout to accept craigslist post")
)

func init() {
	application.Configure()
}

func main() {
	foreachfile.Execute(*craigslistPostLinkDirectoryRootFlag, onFile)
}

func onFile(filePath string) {
	craigslist.Delay()

	url, err := os.ReadFile(filePath)
	logging.Panic(err)

	t, err := time.ParseDuration(*craigslistPostLinkAcceptTimeoutFlag)
	logging.Panic(err)

	tctx, tcancel := context.WithTimeout(application.Context, t)
	defer tcancel()

	craigslist.Accept(tctx, string(url))
}
