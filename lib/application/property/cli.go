package property

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/sflags/gen/gflag"
)

func LoadCli(config interface{}) {
	logging.Panic(gflag.ParseToDef(config))
}
