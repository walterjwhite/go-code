package google

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"google.golang.org/api/option"
)

type pubSubClient interface {
	Publisher(topicNameOrID string) *pubsub.Publisher
	Subscriber(nameOrID string) *pubsub.Subscriber
	Close() error
}

type Conf struct {
	CredentialsFile string
	ProjectId       string

	EncryptionKeyFile string
	Compress          bool
	Serialize         bool

	encryptor encryption.Encryptor

	ctx    context.Context
	cancel context.CancelFunc

	client pubSubClient
}

func (c *Conf) Init(pctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(pctx)

	if len(c.CredentialsFile) == 0 {
		return fmt.Errorf("credentials file path is empty: must be configured")
	}

	if _, err := os.Stat(c.CredentialsFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("credentials file not found")
		}
		return fmt.Errorf("credentials file validation failed")
	}

	realClient, err := pubsub.NewClient(c.ctx, c.ProjectId, option.WithAuthCredentialsFile(option.ServiceAccount, c.CredentialsFile))
	if err != nil {
		return fmt.Errorf("failed to create pubsub client: %w", err)
	}

	if len(c.EncryptionKeyFile) > 0 {
		encryptor, err := aes.NewAESFromFile(c.EncryptionKeyFile)
		if err != nil {
			return fmt.Errorf("failed to initialize encryption: %w", err)
		}

		c.encryptor = encryptor
	}

	c.client = realClient
	return nil
}

func (c *Conf) String() string {
	return fmt.Sprintf("Conf: {ProjectId: %s}", c.ProjectId)
}

func (c *Conf) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}
