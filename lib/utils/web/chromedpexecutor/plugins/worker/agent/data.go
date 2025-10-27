package agent

import (
	"fmt"
	"github.com/walterjwhite/go-code/lib/utils/ui/windows"
)

type Conf struct {
	BrowserName string
	Url         string

	WindowsConf *windows.WindowsConf

	QuestionFile string

	questions   []string
	iteration   int
	contextuals []interface{}
}

func (c *Conf) String() string {
	return fmt.Sprintf("agent.Conf{%s, %s}", c.BrowserName, c.Url)
}
