package property

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"reflect"
)

var (
	e *encryption.EncryptionConfiguration
)

// support strings only
func handleEncryptedProperties(config Configuration) {
	log.Info().Msg("Handling encrypted properties (if any)")
	if config.EncryptedFields() != nil {
		setupEncryption()

		val := reflect.ValueOf(config).Elem()

		for _, fieldName := range config.EncryptedFields() {
			setFieldValue(config, val, fieldName)
		}
	}
}

func setupEncryption() {
	if e == nil {
		log.Info().Msg("Setting up encryption instance")

		e = encryption.New()
	}
}

func setFieldValue(config Configuration, value reflect.Value, fieldName string) {
	log.Info().Msgf("decrypting: %v: %v / %v", value, fieldName, config)

	f := value.FieldByName(fieldName)
	data, err := base64.StdEncoding.DecodeString(f.String())
	logging.Panic(err)

	decrypted := e.Decrypt(data)

	f.SetString(string(decrypted))
}
