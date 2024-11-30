package secrets

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os/exec"
)

func Get(secretName string) string {
	log.Debug().Msgf("secretName: %v", secretName)

	// this fails to work because the cmd is detecting non-interactive and automatically printing debug output and formatting the output for human readability
	cmd := exec.Command("secrets", "get", "-out=stdout", secretName)
	cmd.Env = append(cmd.Environ(), "_FORCE_INTERACTIVE=1")

	out, err := cmd.Output()
	if len(out) == 0 {
		logging.Panic(err)
	}

	log.Debug().Msgf("output: %v", out)
	return string(out[:])
}
