package property

import (
	"flag"
	"github.com/octago/sflags/gen/gflag"
	"github.com/walterjwhite/go-application/libraries/application/logging"
)

func (c *Configuration) LoadCli(config interface{}) {
	logging.Panic(gflag.ParseToDef(config))
	flag.Parse()
}
