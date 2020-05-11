package property

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/typename"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"os"
	"path/filepath"
)

const (
	defaultsPath = "~/.defaults"
)

func (c *Configuration) LoadFile(config interface{}, prefix string) {
	filename := c.getFile(config, prefix)

	finfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Warn().Msgf("%v does not exist", filename)
		return
	}

	if finfo.IsDir() {
		return
	}

	yamlhelper.Read(filename, config)
}

func (c *Configuration) getFile(config interface{}, prefix string) string {
	if len(c.Path) == 0 {
		path, err := homedir.Expand(defaultsPath)
		logging.Panic(err)

		return filepath.Join(path, prefix, typename.Get(config)+".yaml")
	}

	return c.Path
}
