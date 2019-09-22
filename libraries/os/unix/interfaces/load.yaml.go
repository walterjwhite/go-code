package interfaces

import (
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type InterfaceConfiguration struct {
	Interfaces []string
}

var interfaceConfiguration InterfaceConfiguration

func init() {
	yamlhelper.Read("interfaces.yaml", &interfaceConfiguration)

	prepareInterfaces()
}

func prepareInterfaces() {
	//interfaces = [len(interfaceConfiguration.Interfaces)]Interface
	interfaces = make([]Interface, 0)

	for i := 0; i < len(interfaceConfiguration.Interfaces); i++ {
		interfaces = append(interfaces, Interface{interfaceConfiguration.Interfaces[i], uninitialized, false})
	}
}
