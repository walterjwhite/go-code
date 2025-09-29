package citrix

import (
	"errors"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/plugins/mouse_wiggle"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/plugins/noop"
	"time"
)

type WorkerType int

const (
	movementWaitTime   = 3 * time.Minute
	timeBetweenActions = 3 * time.Second
)

const (
	MouseWiggler WorkerType = iota
	NOOP
)

func (w WorkerType) String() string {
	return [...]string{"MouseWiggler", "NOOP"}[w]
}

func (i *Instance) initializeWorker() {
	switch i.WorkerType {
	case MouseWiggler:
		i.Worker = mouse_wiggle.New(movementWaitTime, timeBetweenActions)
	case NOOP:
		i.Worker = &noop.State{}
	default:
		logging.Panic(errors.New("WorkerType unspecified"))
	}
}
