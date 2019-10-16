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
		log.Info().Msgf("Application Version: %v", ApplicationVersion)
		log.Info().Msgf("Built on: %v", BuildDate)
		log.Info().Msgf("OSArchitecture: %v", OSArchitecture)
		log.Info().Msgf("GoVersion: %v", GoVersion)
	} else {
		log.Warn().Msg("Application was not built properly to log application version, build date, etc., check compilation")
	}
}

func isConfigured() bool {
	return len(ApplicationVersion) > 0 && len(BuildDate) > 0
}

func GetApplicationId() string {
	return ApplicationName + "." + ApplicationVersion + "." + SCMId
}
