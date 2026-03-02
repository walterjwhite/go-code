package property

import (
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

type RequiredFields interface {
	RequiredFields() []string
}

func validateRequiredFields(config any) {
	requiredFieldsCfg, ok := config.(RequiredFields)
	if !ok {
		log.Debug().Msgf("%v does not implement RequiredFields, not validating", typename.Get(config))
		return
	}

	fields := requiredFieldsCfg.RequiredFields()
	if fields != nil {
		log.Debug().Msgf("Handling required properties: %v", config)

		rv := reflect.ValueOf(config)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			log.Warn().Msgf("validateRequiredFields expects a pointer to a struct; got: %v", rv.Kind())
			return
		}

		val := rv.Elem()
		if val.Kind() != reflect.Struct {
			log.Warn().Msgf("validateRequiredFields expects a pointer to a struct; got: %v", val.Kind())
			return
		}

		var errCount = 0
		for _, fieldName := range fields {
			if !isValid(requiredFieldsCfg, val, fieldName) {
				errCount += 1
			}
		}

		if errCount > 0 {
			logging.Error(fmt.Errorf("see above warnings regarding required fields"), "validateRequiredFields")
		}
	}
}

func isValid(config RequiredFields, value reflect.Value, fieldName string) bool {
	f := getField(value, fieldName)

	if !f.IsValid() {
		log.Warn().Msgf("field %s not found or invalid", fieldName)
		return false
	}

	if f.Kind() != reflect.String {
		log.Warn().Msgf("field %s is not a string, skipping", fieldName)
		return false
	}

	if !f.CanSet() {
		log.Warn().Msgf("field %s cannot be set (unexported?), skipping", fieldName)
		return false
	}

	fieldValue := f.String()
	log.Debug().Msgf("fieldValue: %v / %v / %v", f, fieldValue, fieldName)
	if len(fieldValue) == 0 {
		log.Warn().Msgf("field %s is empty", fieldName)
		return false
	}

	return true
}
