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

func LoadSecrets(config any) {
	secretPropertyConfiguration, ok := config.(SecretPropertyConfiguration)
	if !ok {
		log.Debug().Msgf("%v does not implement SecretPropertyConfiguration, not decrypting", typename.Get(config))
		return
	}

	fields := secretPropertyConfiguration.SecretFields()
	if fields != nil {
		log.Debug().Msgf("Handling encrypted properties: %v", config)

		rv := reflect.ValueOf(config)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
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

var getSecretFunc = secrets.Get

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
	log.Debug().Msgf("processing field: %s", fieldName)
	if len(fieldValue) > 0 && strings.HasPrefix(fieldValue, "secret://") {
		secretName := strings.TrimPrefix(fieldValue, "secret://")
		if strings.TrimSpace(secretName) == "" {
			log.Warn().Msgf("field %s has secret:// prefix but no secret name, skipping", fieldName)
			return
		}

		secretValue := getSecretFunc(secretName)

		f.SetString(secretValue)
		log.Debug().Msgf("field %s decrypted successfully", fieldName)
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

	return getFieldRecurseWithDepth(field, fieldNamePath[1:], 0)
}

func getFieldRecurseWithDepth(value reflect.Value, fieldNamePath []string, depth int) reflect.Value {
	const maxRecursionDepth = 100
	if depth > maxRecursionDepth {
		return reflect.Value{}
	}

	field := reflect.Indirect(value).FieldByName(fieldNamePath[0])

	if len(fieldNamePath) == 1 {
		return field
	}

	return getFieldRecurseWithDepth(field, fieldNamePath[1:], depth+1)
}

func Decrypt(secretKey string) string {
	return getSecretFunc(secretKey)
}
