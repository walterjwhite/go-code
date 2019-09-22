package server

import (
	"github.com/walterjwhite/go-application/libraries/os/linux/interfaces"
	"log"
)

var interfacesData []interfaces.Interface

type InterfaceServer []interfaces.Interface

func (s *InterfaceServer) Interfaces(args *Args, response *[]interfaces.Interface) error {
	*response = interfacesData
	return nil
}

func RefreshInterfaces() {
	interfacesData = interfaces.Data()
	log.Printf("interfaces: %v\n", interfacesData)
}
