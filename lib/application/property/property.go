package property

import (
	"flag"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)


var (
	pathPrefixFlag = flag.String("config-prefix-path", "", "property prefix, ie. if user specifies web/gmail.com/username with prefix of testing, resulting property would be testing/web/gmail.com/username")
)

func Load(config interface{}) {
	flag.Parse()

	LoadFile(config)

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

}

func GetConfigurationDirectory(qualifiers ...string) string {
	configurationDirectory, err := homedir.Expand("~/.config")
	logging.Panic(err)

	e := []string{configurationDirectory, "walterjwhite"}
	e = append(e, qualifiers...)
	return filepath.Join(e...)
}

func GetConfigurationFile(data interface{}, qualifiers ...string) string {
	return filepath.Join(GetConfigurationDirectory(qualifiers...), typename.Get(data)+".yaml")
}
