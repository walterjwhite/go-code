package monitor

import (
	"github.com/walterjwhite/go-code/lib/time/after"
	"github.com/walterjwhite/go-code/lib/time/periodic"
)

type wrapPeriodic struct {
	Function func()
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
	wrappedPeriodic := &wrapPeriodic{Function: a.Monitor.Execute}

	go periodic.Now(a.Session.Context, periodic.GetInterval(a.Interval), wrappedPeriodic.wrap)
}

func (w *wrapPeriodic) wrap() error {
	w.Function()

	return nil
}

func (s *Session) scheduleNoActivityAlert() {
	if s.NoActivity.After != nil {
		s.NoActivity.After.Cancel()
	}

	s.NoActivity.After = after.After(s.Context, periodic.GetInterval(s.NoActivity.Interval), s.NoActivityAlert)
}
