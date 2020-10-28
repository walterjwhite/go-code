package monitor

import (
	"github.com/walterjwhite/go/lib/io/yaml"
)

func read(configurationFile string, c *Session) {
	yaml.Read(configurationFile, c)
}
