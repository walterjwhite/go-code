package monitor

import (
	"github.com/walterjwhite/go-application/libraries/after"
	"github.com/walterjwhite/go-application/libraries/periodic"
)

type wrapPeriodic struct {
	Function fn()
}

func (a *Action) schedule() {
	if a.isLongRunning() {
		a.invokeLongRunning()
		return
	}

	a.invokePeriodic()
}

func (a *Action) invokeLongRunning() {
	go a.Monitor.Execute()
}

func (a *Action) isLongRunning() bool {
	return len(a.Interval) == 0
}

func (a *Action) invokePeriodic() {
	go periodic.Periodic(a.Session.Context, periodic.GetInterval(a.Interval), &wrapPeriodic{Function: a.Monitor.Execute}.wrap)
}

func (w *wrapPeriodic) wrap() error {
	w.Function()

	return nil
}

func (s *Session) scheduleNoActivityAlert() {
	if s.NoActivity.Timer != nil {
		s.NoActivity.Timer.Stop()
	}

	s.NoActivity.Timer = after.After(s.Context, periodic.GetInterval(s.NoActivity.Interval), s.NoActivityAlert)
}
