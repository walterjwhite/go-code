package property

import (
	"flag"
)

type Configuration struct {
	Path string
}

// TODO: allow other sources (REDIS, etcd, etc.)
// TODO: allow writing properties
type ConfigurationReader interface {
	Load(config interface{})
}

var (
	prefixFlag = flag.String("config-prefix", "", "property prefix, ie. if user specifies web/gmail.com/username with prefix of testing, resulting property would be testing/web/gmail.com/username")
)

func (c *Configuration) Load(config interface{}) {
	c.LoadFile(config)

	c.LoadEnv(config)
	c.LoadCli(config)
	c.LoadSecrets(config)
}

func Load(config interface{}) {
	c := &Configuration{Path: *prefixFlag}
	c.Load(config)
}
