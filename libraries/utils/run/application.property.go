package run

import (
	"github.com/walterjwhite/go-application/libraries/application/property"
)

type ApplicationConfigurer interface {
	Load(a *Application, prefix string)
}

type config struct{}

var (
	Configurer ApplicationConfigurer
)

func New(applications ...string) *Instance {
	setup()

	i := &Instance{}

	i.Applications = make([]*Application, len(applications))
	for index, application := range applications {
		var a *Application

		Configurer.Load(a, application)
		i.Applications[index] = a
	}

	return i
}

func (c *config) Load(a *Application, prefix string) {
	property.Load(a, prefix)
}

func setup() {
	if Configurer == nil {
		Configurer = &config{}
	}
}
