package secrets

import (
	"log"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func setupRepository() {
	if isRepositorySetup() {
		return
	}

	cmd := exec.Command("git", "clone", SecretsConfigurationInstance.repositoryRemoteUri, SecretsConfigurationInstance.repositoryPath)
	stdoutStderr, err := cmd.CombinedOutput()
	logging.Panic(err)

	log.Printf("Setup secrets project: %s\n", stdoutStderr)
}

func isRepositorySetup() bool {
	if _, err := os.Stat(SecretsConfigurationInstance.repositoryPath); os.IsNotExist(err) {
		log.Printf("Secrets !: %v\n", SecretsConfigurationInstance.repositoryPath)
		return false
	}

	log.Printf("Secrets: %v\n", SecretsConfigurationInstance.repositoryPath)
	return true
}
