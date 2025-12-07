package citrix

import (
	"context"
	"errors"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/agent"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/mouse_driver"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker/noop"
)

func (i *Instance) PostLoad(ctx context.Context) error {
	switch i.WorkerType {
	case worker.MouseDriver:
		i.Worker = &mouse_driver.Conf{}
	case worker.Agent:
		i.Worker = &agent.Conf{}
	case worker.NOOP:
		i.Worker = &noop.State{}
	default:
		logging.Panic(errors.New("WorkerType unspecified"))
	}

	application.Load(i.Worker)

	return nil
}
