package codesearch

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

func (s *SearchInstance) file(name string) {
	f, err := os.Open(name)
	if err != nil {
		logging.Warn(err, false, fmt.Sprintf("Unable to search file (%v), is index stale?", name))
		return
	}

	defer f.Close()
	s.reader(f, name)
}
