package worker

import (
	"context"
	"sync"
	"testing"
	"time"
)

type mockWorker struct {
	mu              sync.Mutex
	workCount       int
	shortBreakCount int
	longBreakCount  int
	lunchCount      int
	stopCount       int
}

func (m *mockWorker) String() string { return "mockWorker" }
func (m *mockWorker) Work() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workCount++
}
func (m *mockWorker) ShortBreak() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shortBreakCount++
}
func (m *mockWorker) LongBreak() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.longBreakCount++
}
func (m *mockWorker) Lunch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lunchCount++
}
func (m *mockWorker) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopCount++
}

func (m *mockWorker) getWorkCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.workCount
}

func TestWorker_Run(t *testing.T) {
	mw := &mockWorker{}
	conf := &Conf{
		WorkDuration:       10 * time.Millisecond,
		WorkTickInterval:   5 * time.Millisecond,
		ShortBreakDuration: 5 * time.Millisecond,
		LongBreakDuration:  10 * time.Millisecond,
		LunchStartHour:     time.Now().Hour(),
		LunchDuration:      5 * time.Millisecond,
		StartHour:          0,
		EndHour:            time.Now().Hour() + 1,
	}
	conf.WithWorker(mw)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	go conf.Run(ctx)

	time.Sleep(60 * time.Millisecond)

	if mw.getWorkCount() == 0 {
		t.Error("expected workCount > 0")
	}
}

func TestWorker_Concurrency(t *testing.T) {
	mw1 := &mockWorker{}
	mw2 := &mockWorker{}
	conf := &Conf{
		WorkDuration:       10 * time.Millisecond,
		WorkTickInterval:   5 * time.Millisecond,
		ShortBreakDuration: 5 * time.Millisecond,
		LongBreakDuration:  10 * time.Millisecond,
		StartHour:          0,
		EndHour:            time.Now().Hour() + 1,
	}

	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	conf.WithWorker(mw1)
	go conf.Run(ctx1)

	time.Sleep(20 * time.Millisecond)

	count1 := mw1.getWorkCount()
	if count1 == 0 {
		t.Error("mw1 should have worked by now")
	}

	conf.WithWorker(mw2)
	go conf.Run(ctx2)

	time.Sleep(20 * time.Millisecond)


	cancel1()
	time.Sleep(20 * time.Millisecond)

	count1AfterCancel1 := mw1.getWorkCount()
	time.Sleep(50 * time.Millisecond)
	if mw1.getWorkCount() > count1AfterCancel1 {
		t.Errorf("mw1: work count increased after cancel1: %d -> %d (G1 should have used ctx1)", count1AfterCancel1, mw1.getWorkCount())
	}

	count2 := mw2.getWorkCount()
	time.Sleep(50 * time.Millisecond)
	if mw2.getWorkCount() <= count2 {
		t.Error("mw2 should be increasing")
	}

	cancel2()
	time.Sleep(20 * time.Millisecond)
}
