package secrets

import (
	"errors"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/utils/foreachfile"
)

func Find(callback func(filePath string), patterns ...string) {
	if len(patterns) == 0 {
		logging.Panic(errors.New("At least 1 pattern is required."))
	}

	initialize()

	patterns = append(patterns, "/value")

	foreachfile.Execute(SecretsConfigurationInstance.RepositoryPath, callback, patterns...)
}