package secrets

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const DateTimeLayout = "2006/01/02 15:04:05"

func Encrypt(name *string, message *string, data []byte) {
	log.Printf("processing secret: %v\n", *name)

	initialize()

	secretPath := getSecretPath(name)
	secretValuePath := filepath.Join(secretPath, "value")

	encrypted := SecretsConfigurationInstance.EncryptionConfiguration.Encrypt(data)
	logging.Panic(ioutil.WriteFile(getAbsolute(secretValuePath), encrypted, 0644))

	log.Debug().Msgf("Stored secret in %v (%v)", secretValuePath, len(data))

	putLastUpdated(secretPath)
	commit(secretPath, message)
}

func getSecretPath(name *string) string {
	secretPath := filepath.Join(SecretsConfigurationInstance.RepositoryPath, *name)
	logging.Panic(os.MkdirAll(secretPath, 0755))

	return secretPath
}

func putLastUpdated(secretPath string) {
	secretLastUpdatedPath := filepath.Join(secretPath, "last-updated")

	f, err := os.Create(secretLastUpdatedPath)
	logging.Panic(err)

	defer f.Close()

	lastUpdated := getDateTimeLastUpdated()
	_, err = f.Write(lastUpdated)
	logging.Panic(err)

	log.Debug().Msgf("Stored last updated in %v (%v)", secretLastUpdatedPath, lastUpdated)
}

func getDateTimeLastUpdated() []byte {
	t := time.Now()
	formattedDateTime := t.Format(DateTimeLayout)

	return []byte(formattedDateTime)
}

func commit(secretPath string, message *string) {
	cmd := exec.Command("git", "add", secretPath)
	cmd.Dir = SecretsConfigurationInstance.RepositoryPath

	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("%s\n", stdoutStderr)

	logging.Panic(err)

	cmd = exec.Command("git", "commit", secretPath, "-m", *message)
	cmd.Dir = SecretsConfigurationInstance.RepositoryPath
	stdoutStderr, err = cmd.CombinedOutput()
	log.Debug().Msgf("%s", stdoutStderr)

	logging.Panic(err)

	cmd = exec.Command("git", "push")
	cmd.Dir = SecretsConfigurationInstance.RepositoryPath

	stdoutStderr, err = cmd.CombinedOutput()
	logging.Panic(err)

	log.Debug().Msgf("Added secret:\n%s", stdoutStderr)
}
