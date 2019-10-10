package identifier

import (
	"github.com/rs/zerolog/log"
)

// *MUST* be injected at compile time
var (
	ApplicationName, ApplicationVersion, SCMId, BuildDate, GoVersion, OSArchitecture string
)

func Log() {
	if isConfigured() {
		log.Info().Msgf("Application Version: %v\n", ApplicationVersion)
		log.Info().Msgf("Built on: %v\n", BuildDate)
		log.Info().Msgf("OSArchitecture: %v\n", OSArchitecture)
		log.Info().Msgf("GoVersion: %v\n", GoVersion)
	}
}

func isConfigured() bool {
	return len(ApplicationVersion) > 0 && len(BuildDate) > 0
}

func GetApplicationId() string {
	return ApplicationName + "." + ApplicationVersion + "." + SCMId
}
