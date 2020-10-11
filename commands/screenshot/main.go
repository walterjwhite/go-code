package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"github.com/walterjwhite/go-application/libraries/utils/screenshot"
)

var (
	filenameFlag = flag.String("f", "", "path to save screenshot")
)

func init() {
	application.Configure()
}

func main() {
	if len(*filenameFlag) == 0 {
		logging.Panic(errors.New("filename is required"))
	}

	i := screenshot.Default(*filenameFlag)

	i.Wait()
}
