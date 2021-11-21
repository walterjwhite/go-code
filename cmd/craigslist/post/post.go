package main

import (
	"flag"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/foreachfile"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/craigslist"

	"github.com/walterjwhite/go-code/lib/io/yaml"
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
