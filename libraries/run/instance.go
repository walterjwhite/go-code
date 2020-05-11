package run

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (i *Instance) Run(ctx context.Context, region string) {
	for index, a := range i.Applications {
		a.Run(ctx, region, index)
	}

	i.waitForAll()
}

func (i *Instance) waitForAll() {
	for _, a := range i.Applications {
		_, err := a.command.Process.Wait()

		// TODO: push event to channel
		logging.Panic(err)
	}
}
