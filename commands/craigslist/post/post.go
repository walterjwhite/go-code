package main

import (
	"flag"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/utils/foreachfile"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor/plugins/craigslist"

	"github.com/walterjwhite/go-application/libraries/io/yaml"
)

var (
	craigslistPostDirectoryRootFlag = flag.String("CraigslistPostDirectoryRoot", "./", "Root directory containing craigslist posts")
)

func init() {
	application.Configure()
}

func main() {
	foreachfile.Execute(*craigslistPostDirectoryRootFlag, onFile, ".yaml")
}

func onFile(filePath string) {
	// attempt to rate-limit posts
	craigslist.Wait()

	p := &craigslist.CraigslistPost{}

	yaml.Read(filePath, p)
	p.Create(application.Context)
}
