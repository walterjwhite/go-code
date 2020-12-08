package movement

import (
	"github.com/walterjwhite/go/lib/external/spot/data"
	"time"
)

func (m *MovementConfiguration) OnNewRecord(old, new *data.Record) {
	if isSuspend(new) {
		m.updateAndSchedule(m.SuspendDurationTimeout)
		return
	}

	if m.hasMoved(old, new) {
		m.updateAndSchedule(m.MovementDurationTimeout)
		return
	}
}

func (m *MovementConfiguration) updateAndSchedule(duration time.Duration) {
	m.mutex.Lock()
	m.count = 0
	m.mutex.Unlock()

	m.schedule(duration)
}
