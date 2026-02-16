package periodic

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var counter int32
	fn := func() error {
		atomic.AddInt32(&counter, 1)
		return nil
	}

	p := Now(ctx, cancel, 10*time.Millisecond, fn)

	time.Sleep(25 * time.Millisecond)
	p.Cancel()

	if atomic.LoadInt32(&counter) < 3 {
		t.Errorf("Expected at least 3 executions, but got %d", atomic.LoadInt32(&counter))
	}
}

func TestAfter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var counter int32
	fn := func() error {
		atomic.AddInt32(&counter, 1)
		return nil
	}

	p := After(ctx, cancel, 10*time.Millisecond, fn)

	time.Sleep(25 * time.Millisecond)
	p.Cancel()

	if atomic.LoadInt32(&counter) < 2 {
		t.Errorf("Expected at least 2 executions, but got %d", atomic.LoadInt32(&counter))
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var counter int32
	fn := func() error {
		atomic.AddInt32(&counter, 1)
		return nil
	}

	p := Now(ctx, cancel, 10*time.Millisecond, fn)

	time.Sleep(5 * time.Millisecond) // let it run once
	p.Cancel()
	time.Sleep(20 * time.Millisecond)

	if atomic.LoadInt32(&counter) > 2 { // should be 1, but with timing issues, can be 2
		t.Errorf("Expected 1 or 2 executions, but got %d", atomic.LoadInt32(&counter))
	}
}
