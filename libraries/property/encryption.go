package property

import (
	"reflect"
)

// support strings only
func handleEncryptedProperties(config Configuration) {
	if config.EncryptedFields() != nil {
		val := reflect.ValueOf(&config).Elem()

		for _, fieldName := range config.EncryptedFields() {
			setFieldValue(config, val, fieldName)
		}
	}
}

func setFieldValue(config Configuration, value reflect.Value, fieldName string) {
	f := value.FieldByName(fieldName)
	decrypted := e.Decrypt(f.Bytes())

	f.SetString(string(decrypted))
}
