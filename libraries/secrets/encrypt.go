package secrets

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"log"
	"os"
	"os/exec"
	"time"
)

const DateTimeLayout = "2006/01/02 15:04:05"

func Encrypt(name *string, message *string, data []byte) {
	log.Printf("processing secret: %v\n", *name)

	setupEncryptionKey()

	secretPath := getSecretPath(name)
	secretValuePath := secretPath + "/value"

	SecretsConfigurationInstance.EncryptionConfiguration.EncryptFile(secretValuePath, data)

	log.Printf("Stored secret in %v (%v)\n", secretValuePath, len(data))

	putLastUpdated(secretPath)
	commit(secretPath, message)
}

func getSecretPath(name *string) string {
	secretPath := SecretsConfigurationInstance.RepositoryPath + "/" + *name
	logging.Panic(os.MkdirAll(secretPath, 0755))

	return secretPath
}

func putLastUpdated(secretPath string) {
	secretLastUpdatedPath := secretPath + "/last-updated"

	f, err := os.Create(secretLastUpdatedPath)
	logging.Panic(err)

	defer f.Close()

	lastUpdated := getDateTimeLastUpdated()
	_, err = f.Write(lastUpdated)
	logging.Panic(err)

	log.Printf("Stored last updated in %v (%v)\n", secretLastUpdatedPath, lastUpdated)
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
	log.Printf("%s\n", stdoutStderr)

	logging.Panic(err)

	cmd = exec.Command("git", "push")
	cmd.Dir = SecretsConfigurationInstance.RepositoryPath

	stdoutStderr, err = cmd.CombinedOutput()
	logging.Panic(err)

	log.Printf("Added secret:\n%s\n", stdoutStderr)
}
