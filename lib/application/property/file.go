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

var (
	propertyConfigurationLocationFlag = flag.String("config-path", "~/.config/walterjwhite", "property config path")
	propertyConfigurationFileFlag     = flag.String("property-file", "", "property file")
)

func LoadFile(config interface{}) {
	LoadFileWithPath(config, getFile(config))
}

func LoadFileWithPath(config interface{}, filename string) {
	finfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Error().Msgf("%v does not exist", filename)
		return
	}

	if finfo.IsDir() {
		log.Warn().Msgf("File is a directory %v", filename)
		return
	}

	log.Warn().Msgf("Reading %v", filename)
	yaml.Read(filename, config)
}

func getFile(config interface{}) string {
	if len(*propertyConfigurationFileFlag) > 0 {
		return *propertyConfigurationFileFlag
	}

	if len(*pathPrefixFlag) == 0 {
		path, err := homedir.Expand(*propertyConfigurationLocationFlag)
		logging.Panic(err)

		return filepath.Join(path, *pathPrefixFlag, typename.Get(config)+".yaml")
	}

	return *pathPrefixFlag
}
