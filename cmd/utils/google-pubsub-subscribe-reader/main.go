package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"
)

type ReadSubscriberConfiguration struct {
	CredentialsFile string
	ProjectId       string

	TopicName        string
	SubscriptionName string

	EncryptionKeyFilename string
}

type Callback struct {
	data string
}

var (
	googleConf = &ReadSubscriberConfiguration{}
	aesConf    = &aes.Configuration{}
)

func init() {
	application.ConfigureWithProperties(googleConf)

	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	c := &Callback{}
	session := google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Subscribe(googleConf.TopicName, googleConf.SubscriptionName, c)
	application.Wait()
}

func (c *Callback) New() any {
	return &c.data
}

func (c *Callback) MessageDeserialized() {
	log.Info().Msgf("callback: %s", c.data)
}

func (c *Callback) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}
