package citrix

import (
	"github.com/rs/zerolog/log"
)

func (i *Instance) initializeWorker() {
	err := i.Worker.Init(i.ctx, i.session.Conf.Headless, i.WindowsConf)
	if err != nil {
		log.Warn().Msgf("%v - Instance.run.Worker.Init - error: %v", i, err)
		return
	}
}
