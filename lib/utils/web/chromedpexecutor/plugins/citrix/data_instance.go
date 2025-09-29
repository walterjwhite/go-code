package citrix

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/prompt/windows"
	"sync/atomic"
)

type Instance struct {
	Index      int
	WorkerType WorkerType

	FullScreen bool

	Actions []string

	Worker CitrixWorker

	active *atomic.Bool

	RequiresTermsAcceptance bool
	locked                  bool

	WindowsConf *windows.WindowsConf

	ctx    context.Context
	cancel context.CancelFunc

	session *Session
}

func (i *Instance) String() string {
	return fmt.Sprintf("instance.%d, %d", i.Index, i.WorkerType)
}
