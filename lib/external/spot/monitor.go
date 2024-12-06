package spot

import (
	"context"

	"flag"

	"github.com/walterjwhite/go-code/lib/external/spot/client"

	"path/filepath"
	"time"
)

const (
	minRefreshInterval = time.Duration(150 * time.Second)
)

var (
	dataPath     = flag.String("spot-log-path", "~/.data/spot", "Directory to log spot tracking data")
	testDataPath = flag.String("test-spot-data-path", "", "Directory to use for test spot tracking data")
)

func (c *Configuration) Monitor(ctx context.Context) {
	c.ctx = ctx
	c.initActions()
	c.initFetcher()

	c.fetch()
}

func (c *Configuration) initFetcher() {
	if len(*testDataPath) == 0 {
		c.feedFetcher = &client.Feed{Id: c.Session.FeedId}
		return
	}

	c.feedFetcher = &client.FileFeed{Filename: filepath.Join(*testDataPath, c.Session.FeedId)}
}
