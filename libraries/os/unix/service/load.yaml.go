package service

import (
	"../yamlhelper"
)

type ServicesConfiguration struct {
	Services []*Service
}

var servicesConfiguration ServicesConfiguration

func loadServiceConfiguration() {
	//servicesConfiguration = &ServicesConfiguration{}
	yamlhelper.Read("services.yaml", &servicesConfiguration)
}
