package property

type Configuration struct {
	Path string
}

// TODO: allow other sources (REDIS, etcd, etc.)
// TODO: allow writing properties
type ConfigurationReader interface {
	Load(config interface{}, prefix string)
}

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
