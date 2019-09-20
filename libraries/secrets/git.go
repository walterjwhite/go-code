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

func (s *SecretsConfiguration) Setup() {
	setup := s.isSetup()
	if setup {
		return
	}

	cmd := exec.Command("git", "clone", s.RepositoryRemoteUri, s.RepositoryPath)
	stdoutStderr, err := cmd.CombinedOutput()
	logging.Panic(err)

	log.Printf("Setup secrets project: %s\n", stdoutStderr)
}

func (s *SecretsConfiguration) isSetup() bool {
	if _, err := os.Stat(s.RepositoryPath); os.IsNotExist(err) {
		log.Printf("Secrets !: %v\n", s.RepositoryPath)
		return false
	}

	log.Printf("Secrets: %v\n", s.RepositoryPath)
	return true
}
