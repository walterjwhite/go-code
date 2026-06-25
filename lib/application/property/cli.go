package property

import (
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/urfave/sflags/gen/gflag"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func LoadCli(config any) {
	log.Warn().Msg("Loading CLI properties")
	logging.Warn(gflag.ParseTo(config, flag.CommandLine), "LoadCli")
	flag.Parse()
}
