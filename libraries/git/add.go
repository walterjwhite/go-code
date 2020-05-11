package git

import (
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (c *WorkTreeConfig) Add(filenames ...string) {
	for _, filename := range filenames {
		_, err := c.W.Add(filename)
		logging.Panic(err)
	}
}
