package learning

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"

	"time"
)

type Course struct {
	Title          string
	Url            string
	DurationString string
	Duration       time.Duration
}

func (c *Course) String() string {
	return fmt.Sprintf("{Title: %s, Url: %s, DurationString: %s, Duration: %v}", c.Title, c.Url, c.DurationString, c.Duration)
}

type Session struct {
	EmailAddress string
	Password     string

	StepTimeout           *time.Duration
	WatchIterationTimeout *time.Duration

	AuthRetryAttempts uint
	AuthRetryDelay    time.Duration

	SearchCriteria string

	Conf *provider.Conf

	ctx    context.Context
	cancel context.CancelFunc
}
