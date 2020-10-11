package property

type Configuration struct {
	Path string
}

// TODO: allow other sources (REDIS, etcd, etc.)
// TODO: allow writing properties
type ConfigurationReader interface {
	Load(config interface{}, prefix string)
}

// var (
// 	prefixFlag  = flag.String("prefix", "", "property prefix, ie. if user specifies web/gmail.com/username with prefix of testing, resulting property would be testing/web/gmail.com/username")
// )

func (c *Configuration) Load(config interface{}, prefix string) {
	c.LoadFile(config, prefix)

	LoadEnv(config, prefix)
	LoadCli(config, prefix)
	LoadEncrypted(config, prefix)
}

func Load(config interface{}, prefix string) {
	c := &Configuration{}
	c.Load(config, prefix)
}
