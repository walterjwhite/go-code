package property

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/security/secrets"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

type SecretPropertyConfiguration interface {
	SecretFields() []string
}

func LoadSecrets(config interface{}) {
	secretPropertyConfiguration, ok := config.(SecretPropertyConfiguration)
	if !ok {
		log.Debug().Msgf("%v does not implement SecretPropertyConfiguration, not decrypting", typename.Get(config))
		return
	}

	if secretPropertyConfiguration.SecretFields() != nil {
		log.Debug().Msgf("Handling encrypted properties: %v", config)
		val := reflect.ValueOf(config).Elem()

		log.Debug().Msgf("val: %v", val)

		for _, fieldName := range secretPropertyConfiguration.SecretFields() {
			log.Debug().Msgf("fieldName: %v", fieldName)

			setFieldValue(secretPropertyConfiguration, val, fieldName)
		}
	}
}

func setFieldValue(config SecretPropertyConfiguration, value reflect.Value, fieldName string) {
	f := getField(value, fieldName)

	fieldValue := f.String()
	log.Debug().Msgf("fieldValue: %v / %v / %v", f, fieldValue, fieldName)
	if len(fieldValue) > 0 {
		decrypted := secrets.Get(fieldValue)

		f.SetString(decrypted)
	}
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

func Decrypt(secretKey string) string {
	return secrets.Get(secretKey)
}
