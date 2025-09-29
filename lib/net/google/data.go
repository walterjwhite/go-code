package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"google.golang.org/api/option"
)

type Conf struct {
	CredentialsFile string
	ProjectId       string

	EncryptionKeyFile string
	Compress          bool
	Serialize         bool

	aes *aes.AES

	ctx    context.Context
	Cancel context.CancelFunc

	client *pubsub.Client
}

func (c *Conf) Init(pctx context.Context) {
	c.ctx, c.Cancel = context.WithCancel(pctx)

	client, err := pubsub.NewClient(c.ctx, c.ProjectId, option.WithCredentialsFile(c.CredentialsFile))
	logging.Panic(err)

	if len(c.EncryptionKeyFile) > 0 {
		aes, err := aes.FromFile(c.EncryptionKeyFile)
		logging.Panic(err)

		c.aes = aes
	}

	c.client = client
}

func (c *Conf) String() string {
	return fmt.Sprintf("Conf: {CredentialsFile: %s, ProjectId: %s}", c.CredentialsFile, c.ProjectId)
}
