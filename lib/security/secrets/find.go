package secrets

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/foreachfile"
)

func Find(callback func(filePath string), patterns ...string) {
	if len(patterns) == 0 {
		logging.Panic(errors.New("at least 1 pattern is required"))
	}

	initialize()

	foreachfile.Execute(SecretsConfigurationInstance.RepositoryPath, callback, patterns...)
}
