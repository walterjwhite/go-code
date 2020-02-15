package property

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/secrets"
	"reflect"
)

type encryptionReader struct{}

type SecretPropertyConfiguration interface {
	EncryptedFields() []string
}

func (e *encryptionReader) Load(config interface{}, prefix string) {
	//if ! config instanceof SecretPropertyConfiguration {
	secretPropertyConfiguration, ok := config.(SecretPropertyConfiguration)
	if !ok {
		return
	}

	if secretPropertyConfiguration.EncryptedFields() != nil {
		log.Debug().Msg("Handling encrypted properties")
		val := reflect.ValueOf(config).Elem()

		for _, fieldName := range secretPropertyConfiguration.EncryptedFields() {
			setFieldValue(secretPropertyConfiguration, val, fieldName)
		}
	}
}

func setFieldValue(config SecretPropertyConfiguration, value reflect.Value, fieldName string) {
	f := value.FieldByName(fieldName)

	if len(f.String()) > 0 {
		decrypted := secrets.Decrypt(f.String())

		f.SetString(decrypted)
	}
}
