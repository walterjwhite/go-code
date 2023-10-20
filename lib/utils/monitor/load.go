package monitor

import (
	"github.com/walterjwhite/go-code/lib/io/yaml"
)

func read(configurationFile string, c *Session) {
	yaml.Read(configurationFile, c)
}
