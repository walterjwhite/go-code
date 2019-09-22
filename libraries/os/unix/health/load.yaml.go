package health

import (
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

const (
	systemNetworkConfigurationFile = "/etc/system-status/network.yaml"
)

type NetworkTest struct {
	Type   string
	Target string
}

type NetworkConfiguration struct {
	Tests []NetworkTest
}

var networkConfiguration NetworkConfiguration

func init() {
	yamlhelper.Read(systemNetworkConfigurationFile, &networkConfiguration)
}
