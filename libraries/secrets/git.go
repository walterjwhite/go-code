package secrets

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func setupRepository() {
	if isRepositorySetup() {
		return
	}

	cmd := exec.Command("git", "clone", SecretsConfigurationInstance.RepositoryRemoteUri, SecretsConfigurationInstance.RepositoryPath)
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
