package ssh

import (
	"fmt"

	"os"

	"github.com/walterjwhite/go-code/lib/application/logging"
	s "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func initDirect() ([]byte, error) {
	privateKey := read("Enter Path to Private Key", false)

	key, err := os.ReadFile(*privateKey)
	if err != nil {
		return nil, err
	}

	signer, err := getSigner(key)
	if err != nil {
		return nil, err
	}

	return signer.PublicKey().Marshal(), nil
}

func getSigner(key []byte) (s.Signer, error) {
	signer, err := s.ParsePrivateKey(key)
	if err != nil {
		privateKeyPassPhrase := read("Enter Private Key PassPhrase", true)
		return s.ParsePrivateKeyWithPassphrase(key, []byte(*privateKeyPassPhrase))
	}

	return signer, err
}

func read(message string, hide bool) *string {
	fmt.Println(message)

	if hide {
		value, err := term.ReadPassword(int(os.Stdin.Fd()))
		logging.Panic(err)

		data := string(value)
		return &data
	}

	var value string
	fmt.Scanln(&value)

	return &value
}
