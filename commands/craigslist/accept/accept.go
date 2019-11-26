package main

// accepts the email to "publish" the post
import (
	"context"
	"flag"
	"io/ioutil"
	"time"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/craigslist"
	"github.com/walterjwhite/go-application/libraries/foreachfile"

	"github.com/walterjwhite/go-application/libraries/logging"
)

// approve post
//*[@id="new-edit"]/div/div[4]/div[1]/button
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
	// attempt to rate-limit posts
	craigslist.Wait()

	url, err := ioutil.ReadFile(filePath)
	logging.Panic(err)

	t, err := time.ParseDuration(*craigslistPostLinkAcceptTimeoutFlag)
	logging.Panic(err)

	tctx, tcancel := context.WithTimeout(application.Context, t)
	defer tcancel()

	craigslist.Accept(tctx, string(url))
}
