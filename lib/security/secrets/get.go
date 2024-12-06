package secrets

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os/exec"
)

func Get(secretName string) string {
	log.Debug().Msgf("secretName: %v", secretName)

	cmd := exec.Command("secrets", "get", "-out=stdout", secretName)
	cmd.Env = append(cmd.Environ(), "_FORCE_INTERACTIVE=1")

	out, err := cmd.Output()
	if len(out) == 0 {
		logging.Panic(err)
	}

	log.Debug().Msgf("output: %v", out)
	return string(out[:])
}
