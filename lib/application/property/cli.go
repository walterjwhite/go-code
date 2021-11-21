package property

import (
	"github.com/walterjwhite/sflags/gen/gflag"
	//"github.com/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config interface{}) {
	logging.Panic(gflag.ParseToDef(config))
}
