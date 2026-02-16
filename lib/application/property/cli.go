package property

import (
	"flag"
	"os"

	"github.com/urfave/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config interface{}) {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)
	logging.Warn(gflag.ParseTo(config, fs), "LoadCli")
	logging.Warn(fs.Parse(os.Args[1:]), "fs.Parse")
}
