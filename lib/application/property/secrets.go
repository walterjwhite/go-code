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

	fields := secretPropertyConfiguration.SecretFields()
	if fields != nil {
		log.Debug().Msgf("Handling encrypted properties: %v", config)

		rv := reflect.ValueOf(config)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			log.Warn().Msgf("LoadSecrets expects a pointer to a struct; got: %v", rv.Kind())
			return
		}

		val := rv.Elem()
		if val.Kind() != reflect.Struct {
			log.Warn().Msgf("LoadSecrets expects a pointer to a struct; got: %v", val.Kind())
			return
		}

		for _, fieldName := range fields {
			setFieldValue(secretPropertyConfiguration, val, fieldName)
		}
	}
}

func setFieldValue(config SecretPropertyConfiguration, value reflect.Value, fieldName string) {
	f := getField(value, fieldName)

	if !f.IsValid() {
		log.Warn().Msgf("field %s not found or invalid", fieldName)
		return
	}

	if f.Kind() != reflect.String {
		log.Warn().Msgf("field %s is not a string, skipping", fieldName)
		return
	}

	if !f.CanSet() {
		log.Warn().Msgf("field %s cannot be set (unexported?), skipping", fieldName)
		return
	}

	fieldValue := f.String()
	log.Debug().Msgf("fieldValue: %v / %v / %v", f, fieldValue, fieldName)
	if len(fieldValue) > 0 && strings.HasPrefix(fieldValue, "secret://") {
		secretName := strings.TrimPrefix(fieldValue, "secret://")
		if strings.TrimSpace(secretName) == "" {
			log.Warn().Msgf("field %s has secret:// prefix but no secret name, skipping", fieldName)
			return
		}

		secretValue := secrets.Get(secretName)

		f.SetString(secretValue)
		return
	}

	log.Debug().Msgf("field %s does not start with secret://; leaving value unchanged", fieldName)
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
