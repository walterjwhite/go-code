package property

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/security/secrets"
	"reflect"
	"strings"
)

type SecretPropertyConfiguration interface {
	EncryptedFields() []string
}

func LoadEncrypted(config interface{}, prefix string) {
	secretPropertyConfiguration, ok := config.(SecretPropertyConfiguration)
	if !ok {
		log.Warn().Msg("Unable to initialize secrets, unable to process encrypted properties")
		return
	}

	if secretPropertyConfiguration.EncryptedFields() != nil {
		log.Debug().Msg("Handling encrypted properties")
		val := reflect.ValueOf(config).Elem()

		log.Debug().Msgf("val: %v", val)

		for _, fieldName := range secretPropertyConfiguration.EncryptedFields() {
			log.Debug().Msgf("fieldName: %v %v", prefix, fieldName)

			setFieldValue(secretPropertyConfiguration, val, prefix, fieldName)
		}
	}
}

func setFieldValue(config SecretPropertyConfiguration, value reflect.Value, prefix, fieldName string) {
	f := getField(value, fieldName)

	fieldValue := f.String()
	log.Debug().Msgf("fieldValue: %v / %v / %v", f, fieldValue, fieldName)
	if len(fieldValue) > 0 {
		decrypted := secrets.Decrypt(getFieldPath(prefix, fieldValue))

		f.SetString(decrypted)
	}
}

func getFieldPath(prefix, fieldValue string) string {
	if len(prefix) > 0 {
		return prefix + "/" + fieldValue + "/value"
	}

	return fieldValue + "/value"
}

func getField(value reflect.Value, fieldName string) reflect.Value {
	fieldNamePath := strings.Split(fieldName, ".")

	field := value.FieldByName(fieldNamePath[0])
	if len(fieldNamePath) == 1 {
		return field
	}

	return getFieldRecurse(field, fieldNamePath[1:])
}
func getFieldRecurse(value reflect.Value, fieldNamePath []string) reflect.Value {
	field := reflect.Indirect(value).FieldByName(fieldNamePath[0])

	if len(fieldNamePath) == 1 {
		return field
	}

	return getFieldRecurse(field, fieldNamePath[1:])
}
