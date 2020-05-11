package codesearch

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

func (i *Instance) Reset() {
	// https://github.com/google/codesearch/blob/4fe90b597ae534f90238f82c7b5b1bb6d6d52dff/cmd/cindex/cindex.go#L86
	logging.Panic(os.Remove(i.IndexPath))
}
