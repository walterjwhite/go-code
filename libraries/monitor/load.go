package monitor

import (
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

func read(configurationFile string, c *Session) {
	yamlhelper.Read(configurationFile, c)
}
