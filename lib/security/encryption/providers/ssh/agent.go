package ssh

import (
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"golang.org/x/crypto/ssh/agent"
)

type Conf struct {
	key []byte
}

func initAgent() ([]byte, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err == nil {
		defer logging.Panic(conn.Close())

		client := agent.NewClient(conn)

		keys, err := client.List()
		if err != nil {
			log.Warn().Msg("error listing")
			return nil, err
		}

		if len(keys) != 1 {
			return nil, fmt.Errorf("expecting 1 private key to be available: %v", len(keys))
		}

		return keys[0].Blob, nil
	}

	log.Warn().Msg("error")
	return nil, err
}

func (c *Conf) GetDecryptionKey() []byte {
	return c.key
}

func (c *Conf) GetEncryptionKey() []byte {
	return c.key
}


