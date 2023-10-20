package keep_alive

import (
	"context"
	"time"

	"github.com/walterjwhite/go-code/lib/time/after"
)

type KeepAlive struct {
	activityChannel chan bool
	ctx             context.Context

	timeout     time.Duration
	onKeepAlive func() error

	after   *after.AfterDelay
	onError func(err error)
}

func New(ctx context.Context, timeout time.Duration, fn func() error) *KeepAlive {
	return &KeepAlive{after: after.New(ctx, &timeout, fn), timeout: timeout, onKeepAlive: fn, ctx: ctx}
}

func (k *KeepAlive) WithChannel(activityChannel chan bool) *KeepAlive {
	if k.activityChannel != nil {
		close(k.activityChannel)
	}

	k.activityChannel = activityChannel
	go k.onActivity()

	return k
}

func (k *KeepAlive) onActivity() {
	<-k.activityChannel
	k.Reset()
}

func (k *KeepAlive) WithErrorHandler(fn func(err error)) *KeepAlive {
	k.onError = fn
	return k
}

func (k *KeepAlive) Cancel() {
	k.after.Cancel()
	k.after = nil
}

func (k *KeepAlive) Resume() {
	k.after = after.New(k.ctx, &k.timeout, k.onKeepAlive)
}

func (k *KeepAlive) Reset() {
	k.Cancel()
	k.Resume()
}
