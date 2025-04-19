package worker

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/timeofday"
	"time"
)

type Conf struct {
	StartTime *timeofday.TimeOfDay
	EndTime   *timeofday.TimeOfDay

	LunchStartTime *timeofday.TimeOfDay
	LunchDuration  *time.Duration

	OnDuration    *time.Duration
	OffDuration   *time.Duration
	BreakDuration *time.Duration

	TickInterval *time.Duration

	breakChannel chan *time.Duration
	stopChannel  chan bool
	tickChannel  <-chan time.Time

	Worker Worker

	cycle    int
	hadLunch bool
}

func (c Conf) String() string {
	return fmt.Sprintf("StartTime: %v, EndTime: %v, LunchStartTime: %v, LunchDuration: %v, OnDuration: %v, OffDuration: %v, BreakDuration, %v, TickInterval: %v", c.StartTime, c.EndTime, c.LunchStartTime, c.LunchDuration, c.OnDuration, c.OffDuration, c.BreakDuration, c.TickInterval)
}

type Worker interface {
	Work()

	OnBreak(*time.Duration)

	OnStop()
}

func (c *Conf) Run() {
	log.Info().Msg("worker.Run()")
	if !c.StartTime.SleepUntil() {
		log.Warn().Msg("start time already passed")
	}
	if c.EndTime.Till() < 0 {
		log.Warn().Msg("end time already passed")
		return
	}

	c.breakChannel = make(chan *time.Duration)
	c.stopChannel = make(chan bool)
	c.tickChannel = time.Tick(*c.TickInterval)

	defer close(c.breakChannel)
	defer close(c.stopChannel)

	go c.manageBreaks()
	go c.stop()

	log.Warn().Msg("running first iteration")
	c.Worker.Work()

	for {
		select {
		case <-c.stopChannel:
			log.Warn().Msgf("exiting instance: %v", c)
			c.Worker.OnStop()

			return
		case duration := <-c.breakChannel:
			log.Warn().Msgf("taking a break: %v", *duration)
			c.Worker.OnBreak(duration)

			time.Sleep(*duration)
		case <-c.tickChannel:
			c.Worker.Work()
		}
	}
}

func (c *Conf) WillRun() bool {
	return c.EndTime.Till() > 0
}
