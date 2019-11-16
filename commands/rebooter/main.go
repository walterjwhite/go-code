package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/os/reboot"
	"github.com/walterjwhite/go-application/libraries/property"
	"net/http"
)

type WebRebootConfiguration struct {
	EndpointPath string

	Port    int
	Address string

	RebootConfiguration *reboot.RebootConfiguration
}

func (c *WebRebootConfiguration) HasDefault() bool {
	return false
}

func (c *WebRebootConfiguration) Refreshable() bool {
	return false
}

func (c *WebRebootConfiguration) EncryptedFields() []string {
	return nil
}

var (
	webRebootConfiguration = &WebRebootConfiguration{}
)

func init() {
	fmt.Println("Main.init()")
	application.Configure( /*rebootConfiguration*/ /*TODO: future versions might do this automatically for us*/ )

	log.Info().Msgf("loading")
	property.Load(webRebootConfiguration, "")
	log.Info().Msgf("loaded")
}

func main() {
	log.Info().Msgf("Listening @: %v", webRebootConfiguration.EndpointPath)
	http.HandleFunc(webRebootConfiguration.EndpointPath, handleReboot)

	log.Info().Msgf("Listening @: %v/%v", webRebootConfiguration.Address, webRebootConfiguration.Port)
	logging.Panic(http.ListenAndServe(fmt.Sprintf("%v:%v", webRebootConfiguration.Address, webRebootConfiguration.Port), nil))
}

func handleReboot(w http.ResponseWriter, r *http.Request) {
	var message string
	if webRebootConfiguration.RebootConfiguration.Reboot(application.Context) {
		message = "Rebooting"
	} else {
		message = "NOT rebooting"
	}

	_, err := w.Write([]byte(message))
	logging.Panic(err)
}
