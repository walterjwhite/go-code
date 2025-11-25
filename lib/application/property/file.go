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
	configFilePrefixFlag = flag.String("conf-prefix", "", "If specified, configuration files will be expected to be nested in this directory, ie. ~/.config/walterjwhite/<prefix>/<TypeName>.yaml")
)

func LoadFile(applicationName string, config interface{}) {
	LoadFileWithPath(config, getFile(applicationName, config))
}

func LoadFileWithPath(config interface{}, filename string) {
	finfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Warn().Msgf("%v does not exist", filename)
		return
	}

	if finfo.IsDir() {
		log.Warn().Msgf("File is a directory %v", filename)
		return
	}

	log.Info().Msgf("Reading %v", filename)
	yaml.Read(filename, config)
}

func getFile(applicationName string, config interface{}) string {
	path, err := homedir.Expand(propertyConfigurationLocation)
	logging.Panic(err)

	return filepath.Join(path, applicationName, *configFilePrefixFlag, typename.Get(config)+".yaml")
}
