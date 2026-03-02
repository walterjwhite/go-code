package property

import (
	"flag"

	"github.com/urfave/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config any) {
	fs := flag.CommandLine 
	logging.Warn(gflag.ParseTo(config, fs), "LoadCli")
	flag.Parse()
}
