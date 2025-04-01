package google

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"google.golang.org/api/option"
)

type Conf struct {
	CredentialsFile string
	ProjectId       string
}

type Session struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	client *pubsub.Client

	AesConf           *aes.Configuration
	EnableCompression bool
}

func New(credentialsFile string, projectId string, pctx context.Context) *Session {
	ctx, cancel := context.WithCancel(pctx)

	client, err := pubsub.NewClient(ctx, projectId, option.WithCredentialsFile(credentialsFile))
	logging.Panic(err)

	return &Session{Ctx: ctx, Cancel: cancel, client: client}
}
