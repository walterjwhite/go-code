package virginpulse

import (
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/run"
)

// just do today for starters
// TODO: put this in a script
func (s *Session) RunScript() {
	for {
		s.ChromeDPSession.Execute(run.ParseActions(s.Script...)...)
	}
}
