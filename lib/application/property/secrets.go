package property

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/security/secrets"
	"reflect"
	"strings"
)

type SecretPropertyConfiguration interface {
	SecretFields() []string
}

func LoadSecrets(config interface{}) {
	secretPropertyConfiguration, ok := config.(SecretPropertyConfiguration)
	if !ok {
		log.Warn().Msg("Unable to initialize secrets, unable to process encrypted properties")
		return
	}

	if secretPropertyConfiguration.SecretFields() != nil {
		log.Debug().Msgf("Handling encrypted properties: %v", config)
		val := reflect.ValueOf(config).Elem()

		log.Debug().Msgf("val: %v", val)

		for _, fieldName := range secretPropertyConfiguration.SecretFields() {
			log.Debug().Msgf("fieldName: %v %v", *pathPrefixFlag, fieldName)

			setFieldValue(secretPropertyConfiguration, val, *pathPrefixFlag, fieldName)
		}
	}
}

func setFieldValue(config SecretPropertyConfiguration, value reflect.Value, path, fieldName string) {
	f := getField(value, fieldName)

	fieldValue := f.String()
	log.Debug().Msgf("fieldValue: %v / %v / %v", f, fieldValue, fieldName)
	if len(fieldValue) > 0 {
		decrypted := secrets.Decrypt(getFieldPath(path, fieldValue))

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

func Decrypt(secretKey string) string {
	return secrets.Decrypt(getFieldPath(*pathPrefixFlag, secretKey))
}
