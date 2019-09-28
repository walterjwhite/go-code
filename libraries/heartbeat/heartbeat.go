package heartbeat

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/periodic"
)

type HeartbeatInstance struct {
	Interval time.Duration

	// Function to invoke every interval
	HeartbeatFunction func() error

	// Long-running function
	Function func()
}

func Heartbeat(function func(), heartbeatFunction func() error, interval time.Duration) {
	h := &HeartbeatInstance{Interval: interval, HeartbeatFunction: heartbeatFunction, Function: function}
	h.Call()
}

func (h *HeartbeatInstance) Call() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go periodic.Periodic(ctx, h.Interval, h.HeartbeatFunction)
	h.Function()
}
