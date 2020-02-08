package property

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/typename"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"os"
	"path/filepath"
)

type defaultConfigurationReader struct{}

var (
	basePath = flag.String("PropertyBasePath", "~/.defaults", "PropertyBasePath")
	//filePath = flag.String("PropertyFilePath", "", "PropertyFilePath")
)

func (d *defaultConfigurationReader) Load(config interface{}, prefix string) {
	defaultsBasePath, err := homedir.Expand(*basePath)
	logging.Panic(err)

	f := &fileConfigurationReader{Filename: filepath.Join(defaultsBasePath, prefix, typename.Get(config)+".yaml")}
	f.Load(config, prefix)
}

type fileConfigurationReader struct {
	Filename string
}

func (f *fileConfigurationReader) Load(config interface{}, prefix string) {
	finfo, err := os.Stat(f.Filename)
	if os.IsNotExist(err) {
		log.Warn().Msgf("%v does not exist", f.Filename)
		return
	}

	if finfo.IsDir() {
		return
	}

	yamlhelper.Read(f.Filename, config)
}
