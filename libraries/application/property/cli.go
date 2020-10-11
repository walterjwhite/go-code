package property

import (
	"github.com/octago/sflags/gen/gflag"

	"github.com/walterjwhite/go-application/libraries/application/logging"
)

func LoadCli(config interface{}, prefix string) {
	_, err := gflag.Parse(config)
	logging.Panic(err)
}
