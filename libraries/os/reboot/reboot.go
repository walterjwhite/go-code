package reboot

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/os/shutdown"
	"github.com/walterjwhite/go-application/libraries/os/user/idle"
	"sync"
	"time"
)

type RebootConfiguration struct {
	DryRun  bool
	Timeout time.Duration

	MinIdleTime time.Duration
}

var (
	rebootMutex = &sync.Mutex{}
)

func (c *RebootConfiguration) Reboot(ctx context.Context) bool {
	log.Info().Msg("Request received")

	rebootMutex.Lock()
	defer rebootMutex.Unlock()

	if c.canReboot(ctx) {
		c.doReboot()
		return true
	}

	c.doNotReboot()
	return false
}

func (c *RebootConfiguration) canReboot(ctx context.Context) bool {
	return idle.IdleTime(ctx) > c.MinIdleTime
}

func (c *RebootConfiguration) doReboot() {
	log.Info().Msg("Rebooting")

	shutdownRequest := shutdown.ShutdownRequest{DryRun: c.DryRun, Timeout: c.Timeout, ShutdownAction: shutdown.Reboot}
	shutdownRequest.Execute()
}

func (c *RebootConfiguration) doNotReboot() {
	log.Info().Msg("Not Rebooting")
}
