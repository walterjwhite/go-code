package ssh

import (
	"crypto/sha256"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

type Conf struct {
	agentClient agent.ExtendedAgent
}

var (
	Instance *Conf
)

func init() {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	logging.Panic(err)

	Instance = &Conf{agentClient: agent.NewClient(conn)}
}

func (c *Conf) List() {
	keys, err := c.agentClient.List()
	logging.Panic(err)

	for _, key := range keys {
		log.Info().Msgf("Key: %v", key)
		log.Info().Msgf("Key: %v", string(key.Blob))
	}
}

func (c *Conf) GetDecryptionKey() []byte {
	return getUsableKey(c.getDecryptionKey().Blob)
}

// TODO: select the key to use (flag, configuration?)
func (c *Conf) getDecryptionKey() *agent.Key {
	keys, err := c.agentClient.List()
	logging.Panic(err)

	if len(keys) != 1 {
		logging.Panic(fmt.Errorf("Expecting 1 private key to be available: %v", len(keys)))
	}

	return keys[0]
}

func (c *Conf) GetEncryptionKey() []byte {
	return c.GetDecryptionKey()
}

// TODO: we should probably hash this several times?
func getUsableKey(data []byte) []byte {
	k := sha256.Sum256(data)
	return k[:]
}
