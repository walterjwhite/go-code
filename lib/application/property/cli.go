package property

import (
	"github.com/urfave/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config interface{}) {
	logging.Panic(gflag.ParseToDef(config))
}
