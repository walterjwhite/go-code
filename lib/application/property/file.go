package property

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/yaml"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

const (
	propertyConfigurationLocation = "~/.config/walterjwhite/go"
)

var (
	configFilePrefixFlag = flag.String("conf-prefix", "", "additional sub-directory to help differentiate between configuration")

	getFileFunc = getFile
)

func LoadFile(applicationName string, config interface{}) error {
	return LoadFileWithPath(config, getFileFunc(applicationName, config))
}

func LoadFileWithPath(config interface{}, filename string) error {
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil
	}

	if finfo.IsDir() {
		return nil
	}

	log.Info().Msgf("Reading %v", filename)
	return yaml.Read(filename, config)
}

func getFile(applicationName string, config interface{}) string {
	if len(applicationName) == 0 {
		log.Warn().Msgf("application name is empty: %s", applicationName)
	}

	path, err := homedir.Expand(propertyConfigurationLocation)
	logging.Error(err)

	return filepath.Join(path, applicationName, *configFilePrefixFlag, typename.Get(config)+".yaml")
}
