package application

import (
	"github.com/rs/zerolog/log"
)

// *MUST* be injected at compile time
var (
	ApplicationName, ApplicationVersion, SCMId, BuildDate, GoVersion, OSArchitecture string
)

func logIdentifier() {
	if !isConfigured() {
		log.Warn().Msg("Application was not built properly to log application version, build date, etc., check compilation")
		return
	}

	log.Debug().Msgf("Application: %v", GetApplicationId())
	log.Debug().Msgf("Built on: %v", BuildDate)
	log.Debug().Msgf("OSArchitecture: %v", OSArchitecture)
	log.Debug().Msgf("GoVersion: %v", GoVersion)
}

func isConfigured() bool {
	return len(ApplicationVersion) > 0 && len(BuildDate) > 0
}

func GetApplicationId() string {
	return ApplicationName + "." + ApplicationVersion + "." + SCMId
}
