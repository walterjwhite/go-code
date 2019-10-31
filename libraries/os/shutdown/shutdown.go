package shutdown

import (
	"github.com/rs/zerolog/log"
	"time"
	)
	
	type ShutdownAction int
	
	const (
	Reboot ShutdownAction = 0
	Poweroff ShutdownAction = 1
	)
	
type ShutdownRequest struct {
	DryRun bool
	Timeout time.Duration
	SHutdownAction ShutdownAction
}

type Shutdowner interface {
	log()
	run()
}

func (r *ShutdownRequest) Execute() {
	if r.DryRun {
		r.doDryRUn()
		return
	}
	
	r.run()
}

func (r *ShutdownRequest) doDryRun() {
	log.Warn().Msg("Performing dry run only")
	r.log()
}
