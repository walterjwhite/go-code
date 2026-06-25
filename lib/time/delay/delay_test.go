package delay

import (
	"testing"
	"time"
)

func TestFixedDelay(t *testing.T) {
	d := New(10 * time.Millisecond)
	start := time.Now()
	d.Delay()
	end := time.Now()
	if end.Sub(start) < 10*time.Millisecond {
		t.Errorf("FixedDelay did not sleep for at least 10ms")
	}

	d = New(0)
	start = time.Now()
	d.Delay()
	end = time.Now()
	if end.Sub(start) > 1*time.Millisecond {
		t.Errorf("FixedDelay with 0 duration slept for too long")
	}
}

func TestRandomDelay(t *testing.T) {
	maxDuration := 20 * time.Millisecond
	d := NewRandom(maxDuration)
	start := time.Now()
	d.Delay()
	end := time.Now()
	if end.Sub(start) > maxDuration+5*time.Millisecond {
		t.Errorf("RandomDelay slept for longer than the max duration: %v", end.Sub(start))
	}

	d = NewRandom(0)
	start = time.Now()
	d.Delay()
	end = time.Now()
	if end.Sub(start) > 1*time.Millisecond {
		t.Errorf("RandomDelay with 0 duration slept for too long")
	}
}
