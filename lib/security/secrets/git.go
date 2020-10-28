package secrets

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"os"
	"os/exec"
)

func setupRepository() {
	if isRepositorySetup() {
		return
	}

	target, err := homedir.Expand(SecretsConfigurationInstance.RepositoryPath)
	logging.Panic(err)

	cmd := exec.Command("git", "clone", SecretsConfigurationInstance.RepositoryRemoteUri, target)
	stdoutStderr, err := cmd.CombinedOutput()
	logging.Panic(err)

	log.Debug().Msgf("Setup secrets project: %s", stdoutStderr)
}

func isRepositorySetup() bool {
	if _, err := os.Stat(SecretsConfigurationInstance.RepositoryPath); os.IsNotExist(err) {
		log.Debug().Msgf("Secrets !: %v", SecretsConfigurationInstance.RepositoryPath)
		return false
	}

	log.Debug().Msgf("Secrets: %v", SecretsConfigurationInstance.RepositoryPath)
	return true
}
