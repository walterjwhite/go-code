package property

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/io/yaml"
	"github.com/walterjwhite/go/lib/utils/typename"
	"os"
	"path/filepath"
)

var (
	propertyConfigurationLocationFlag = flag.String("config-path", "~/.config/walterjwhite", "property config path")
)

func (c *Configuration) LoadFile(config interface{}) {
	filename := c.getFile(config)

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

func (c *Configuration) getFile(config interface{}) string {
	if len(c.Path) == 0 {
		path, err := homedir.Expand(*propertyConfigurationLocationFlag)
		logging.Panic(err)

		return filepath.Join(path, c.Path, typename.Get(config)+".yaml")
	}

	return c.Path
}
