package property

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/io/yaml"
	"github.com/walterjwhite/go/lib/utils/typename"
)

var (
	propertyConfigurationLocationFlag = flag.String("config-path", "~/.config/walterjwhite", "property config path")
	propertyConfigurationFileFlag     = flag.String("property-file", "", "property file")
)

func LoadFile(config interface{}) {
	filename := getFile(config)

	finfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Error().Msgf("%v does not exist", filename)
		return
	}

	if finfo.IsDir() {
		return
	}

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
