package property

import (
	"flag"
	"os"
	"strings"

	"github.com/urfave/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config any) {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return
		}
	}

	fs := cloneFlagSet(flag.CommandLine)
	logging.Warn(gflag.ParseTo(config, fs), "LoadCli")
	logging.Warn(fs.Parse(os.Args[1:]), "LoadCli")
}

func cloneFlagSet(src *flag.FlagSet) *flag.FlagSet {
	fs := flag.NewFlagSet(src.Name(), flag.ContinueOnError)
	src.VisitAll(func(f *flag.Flag) {
		fs.Var(f.Value, f.Name, f.Usage)
	})
	return fs
}
