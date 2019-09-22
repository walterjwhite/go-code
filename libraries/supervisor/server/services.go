package server

import (
	"github.com/walterjwhite/go-application/libraries/unix/service"
	"log"
)

var services []service.Service

type ServiceServer []service.Service

func (s *ServiceServer) Services(args *Args, response *[]service.Service) error {
	*response = services
	return nil
}

func RefreshServices() {
	services = service.Data()
	log.Printf("services: %v", services)
}
