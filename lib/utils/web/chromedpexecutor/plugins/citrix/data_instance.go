package citrix

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/ui/windows"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/worker"
	"sync/atomic"
	"time"
)

type Instance struct {
	Index      int
	WorkerType worker.WorkerType
	Lockable bool

	FullScreen          bool
	InitializationDelay time.Duration

	Actions []string

	Worker worker.ChromeDPWorker

	active *atomic.Bool

	RequiresTermsAcceptance bool
	locked                  bool

	WindowsConf *windows.WindowsConf

	ctx    context.Context
	cancel context.CancelFunc

	session *Session
}

func (i *Instance) String() string {
	return fmt.Sprintf("instance.%d, %s", i.Index, i.WorkerType)
}
