package main

import "github.com/walterjwhite/go-code/lib/application/property"

type Config struct {
	Proxy              string `yaml:"proxy" flag:"proxy p" desc:"SOCKS proxy address (host:port)"`
	DB                 string `yaml:"db" flag:"db" desc:"SQLite database path"`
	Service            string `yaml:"service" flag:"service" desc:"Service to determine public IP"`
	HTTPTimeoutSeconds int    `yaml:"http_timeout_seconds" flag:"http-timeout-seconds" desc:"HTTP timeout in seconds"`
}

var _ property.PreLoad = (*Config)(nil)

func (c *Config) PreLoad() {
	if len(c.DB) == 0 {
		c.DB = DefaultDBPath
	}
	if len(c.Service) == 0 {
		c.Service = DefaultServiceURL
	}
	if c.HTTPTimeoutSeconds == 0 {
		c.HTTPTimeoutSeconds = HTTPTimeout
	}
}
