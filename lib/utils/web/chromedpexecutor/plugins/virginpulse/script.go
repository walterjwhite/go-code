package virginpulse

import (
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"
)

// just do today for starters
// TODO: put this in a script
func (s *Session) RunScript() {
	for {
		s.ChromeDPSession.Execute(chromedpexecutor.ParseActions(s.Script...)...)
	}
}
