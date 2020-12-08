package property

import (
	"flag"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/utils/typename"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

// TODO: allow other sources (REDIS, etcd, etc.)
// TODO: allow writing properties
// type ConfigurationReader interface {
// 	Load(config interface{})
// }

var (
	pathPrefixFlag = flag.String("config-prefix-path", "", "property prefix, ie. if user specifies web/gmail.com/username with prefix of testing, resulting property would be testing/web/gmail.com/username")
)

func Load(config interface{}) {
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
