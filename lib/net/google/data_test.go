package google

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf_String(t *testing.T) {
	conf := &Conf{
		CredentialsFile: "test_credentials.json",
		ProjectId:       "test_project",
	}
	expected := "Conf: {CredentialsFile: test_credentials.json, ProjectId: test_project}"
	assert.Equal(t, expected, conf.String())
}

/*
func TestConf_Init_NoEncryption(t *testing.T) {

	conf := &Conf{
		CredentialsFile: "/dev/null", // Use a dummy file path
		ProjectId:       "test_project",
	}

	pctx := context.Background()
	conf.Init(pctx)

	assert.NotNil(t, conf.ctx)
	assert.NotNil(t, conf.cancel)

}
*/

/*
func TestConf_Init_WithEncryption(t *testing.T) {

	conf := &Conf{
		CredentialsFile:   "/dev/null", // Dummy file
		ProjectId:         "test_project",
		EncryptionKeyFile: "non_existent_key.aes", // This will cause logging.Error without mocking
	}

	pctx := context.Background()
	conf.Init(pctx)

	assert.NotNil(t, conf.ctx)
	assert.NotNil(t, conf.cancel)
}
*/

/*
func TestConf_Cancel(t *testing.T) {
	conf := &Conf{}
	pctx := context.Background()
	conf.Init(pctx) // Initialize context and cancel function

	conf.Cancel()

	select {
	case <-conf.ctx.Done():
	default:
		t.Error("Context was not cancelled")
	}

	conf.Cancel()
}
*/
