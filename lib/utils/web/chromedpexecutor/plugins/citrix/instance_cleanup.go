package citrix

import (
	"github.com/rs/zerolog/log"
)

func (i *Instance) cleanup() {
	if i.Worker != nil {
		i.Worker.Cleanup()
	}

	select {
	case <-i.ctx.Done():
		log.Warn().Msgf("%v - instance.cleanup - context done", i)
		return
	default:
		log.Warn().Msgf("%v - instance.cleanup - cancel context", i)
		i.cancel()
	}
}
