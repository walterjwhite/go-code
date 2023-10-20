package ssh

import (
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/agent"
)

type Conf struct {
	key []byte
}

func initAgent() ([]byte, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err == nil {
		defer conn.Close()

		client := agent.NewClient(conn)

		// this was working before, but is not now?
		keys, err := client.List()
		if err != nil {
			log.Warn().Msg("error listing")
			return nil, err
		}

		if len(keys) != 1 {
			return nil, fmt.Errorf("Expecting 1 private key to be available: %v", len(keys))
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

// func (c *Conf) List() {
// 	keys, err := c.agentClient.List()
// 	logging.Panic(err)

// 	for _, key := range keys {
// 		log.Info().Msgf("Key: %v", key)
// 		log.Info().Msgf("Key: %v", string(key.Blob))
// 	}
// }
