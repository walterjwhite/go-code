package secrets

import (
	"log"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type Session struct {
	Path *string
}

func setupRepository() {
	if isRepositorySetup() {
		return
	}

	cmd := exec.Command("git", "clone", SecretsConfigurationInstance.RepositoryRemoteUri, SecretsConfigurationInstance.RepositoryPath)
	stdoutStderr, err := cmd.CombinedOutput()
	logging.Panic(err)

	log.Printf("Setup secrets project: %s\n", stdoutStderr)
}

func isRepositorySetup() bool {
	if _, err := os.Stat(SecretsConfigurationInstance.RepositoryPath); os.IsNotExist(err) {
		log.Printf("Secrets !: %v\n", SecretsConfigurationInstance.RepositoryPath)
		return false
	}

	log.Printf("Secrets: %v\n", SecretsConfigurationInstance.RepositoryPath)
	return true
}
