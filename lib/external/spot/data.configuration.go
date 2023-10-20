package spot

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/action"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/external/spot/writer"
	"github.com/walterjwhite/go-code/lib/external/spot/writer/jsonfile"
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"os"
	"path/filepath"
	"time"
)

type AlertLevel int

const (
	Warning AlertLevel = iota
	Critical
)

type Configuration struct {
	Session *data.Session

	RefreshInterval time.Duration
	Actions         []action.BackgroundAction

	feedFetcher   client.FeedFetcher
	fetchPeriodic *periodic.PeriodicInstance

	writer writer.SpotWriter

	ctx context.Context
}

func New(feedId string) *Configuration {
	c := &Configuration{Session: &data.Session{FeedId: feedId}}

	c.writer = jsonfile.New(c.Session)

	c.initPaths()
	c.initLastRecord()

	return c
}

func (c *Configuration) initPaths() {
	path, err := homedir.Expand(*dataPath)
	logging.Panic(err)

	c.Session.SessionPath = filepath.Join(path, c.Session.FeedId)
	c.Session.DataPath = filepath.Join(c.Session.SessionPath, ".data")

	logging.Panic(os.MkdirAll(path, os.ModePerm))
	logging.Panic(os.MkdirAll(c.Session.DataPath, os.ModePerm))

	if len(*testDataPath) > 0 {
		logging.Panic(os.MkdirAll(*testDataPath, os.ModePerm))
	}
}
