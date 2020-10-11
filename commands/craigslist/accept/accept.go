package main

// accepts the email to "publish" the post
import (
	"context"
	"flag"
	"io/ioutil"
	"time"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/utils/foreachfile"
	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor/plugins/craigslist"

	"github.com/walterjwhite/go-application/libraries/application/logging"
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

// TODO:
// 1. read IMAP messages
// 2. automatically accept posts
// 3. move confirmation emails to another folder
// 4. automatically respond back confirming whether an item exists?
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
