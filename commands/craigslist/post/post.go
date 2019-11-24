package main

import (
	"flag"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/craigslist"
	"github.com/walterjwhite/go-application/libraries/foreachfile"

	"github.com/walterjwhite/go-application/libraries/yamlhelper"
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

	yamlhelper.Read(filePath, p)
	p.Create(application.Context)
}
