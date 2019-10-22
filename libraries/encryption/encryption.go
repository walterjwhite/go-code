package encryption

import (
	"bufio"
	"errors"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

type EncryptionConfiguration struct {
	encryptionKey []byte
}

func New() *EncryptionConfiguration {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		keyBytes := scanner.Bytes()
		keyBytes = append(keyBytes, '\n')

		return &EncryptionConfiguration{encryptionKey: keyBytes}
	}

	logging.Panic(errors.New("No encryption key provided"))
	return nil
}
