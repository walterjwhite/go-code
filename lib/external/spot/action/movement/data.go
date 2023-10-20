package movement

import (
	"context"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"github.com/walterjwhite/go-code/lib/time/after"
	"strings"
	"sync"
	"time"
)

type MovementConfiguration struct {
	Session *data.Session

	MovementTolerance float64

	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int

	AlertAfter int

	MovementDurationTimeout time.Duration

	// once user presses okay, how long should we disable monitoring the user's movement?
	SuspendDurationTimeout time.Duration

	//monitorPeriodic *periodic.PeriodicInstance
	after *after.AfterDelay
	count int

	parentContext context.Context

	mutex *sync.RWMutex
}

func New(s *data.Session) *MovementConfiguration {
	return &MovementConfiguration{Session: s, mutex: &sync.RWMutex{}}
}

func (c *MovementConfiguration) Init(s *data.Session, ctx context.Context) {
	c.parentContext = ctx

	c.schedule(c.getDuration())
}

func (c *MovementConfiguration) getDuration() time.Duration {
	if c.Session.LatestReceivedRecord == nil || strings.Compare(string(client.OK), string(c.Session.LatestReceivedRecord.MessageType)) != 0 {
		return c.MovementDurationTimeout
	}

	return c.SuspendDurationTimeout
}
