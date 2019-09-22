package service

import (
	"../data"
	"../health"
	"log"
)

const (
	RUNNING  = "Running"
	STOPPED  = "Stopped"
	STARTING = "Starting"
)

type Service struct {
	Name      string
	Port      string
	Check     string
	Arguments []string

	Status string
	Health int

	Addresses []string
	//Addresses []Address

	Restart []string
}

type Address struct {
	IP   string
	Port int
}

func Header() *data.Header {
	return &header
}

//var SERVICES []Service

func Data() []Service {
	channel := make(chan *Service)

	refreshAll(channel)
	waitForAllToComplete(channel)

	log.Printf("services (refresh): %v", servicesConfiguration.Services)

	//return servicesConfiguration.Services

	services := make([]Service, 0)

	for i := 0; i < len(servicesConfiguration.Services); i++ {
		services = append(services, *servicesConfiguration.Services[i])
	}

	return services
}

func refreshAll(channel chan *Service) {
	for _, service := range servicesConfiguration.Services {
		log.Printf("refreshAll(%v) %v %v\n", service.Name, &service, service)

		go call(service, channel)
	}
}

func waitForAllToComplete(channel chan *Service) {
	for range servicesConfiguration.Services {
		// result := <- channel
		// services = append(services, result)
		<-channel
	}

	log.Printf("all completed")
}

func call(service *Service, channel chan *Service) {
	log.Printf("call: %v\n", service.Name)

	service.Status = getServiceStatus(service)
	service.Addresses = getServiceListeningAddresses(service)
	service.Health = checkHealth(service)

	log.Printf("Service (service.go): %v %v %v %v\n", service.Name, service.Status, service.Addresses, service.Health)

	// queue a repair if required
	repair(service)

	channel <- service
}

func checkHealth(service *Service) int {
	return health.Health(service.Check, service.Arguments)
}

func init() {
	loadServiceConfiguration()
}
