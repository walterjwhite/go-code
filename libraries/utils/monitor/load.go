package monitor

import (
	"github.com/walterjwhite/go-application/libraries/io/yaml"
)

func read(configurationFile string, c *Session) {
	yaml.Read(configurationFile, c)
}
