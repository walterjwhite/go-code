package identifier

import (
	"log"
)

// *MUST* be injected at compile time
var (
	ApplicationName, ApplicationVersion, SCMId, BuildDate, GoVersion, OSArchitecture string
)

func Log() {
	if isConfigured() {
		log.Printf("Application Version: %v\n", ApplicationVersion)
		log.Printf("Built on: %v\n", BuildDate)
		log.Printf("OSArchitecture: %v\n", OSArchitecture)
		log.Printf("GoVersion: %v\n", GoVersion)
	}
}

func isConfigured() bool {
	return len(ApplicationVersion) > 0 && len(BuildDate) > 0
}

func GetApplicationId() string {
	return ApplicationName + "." + ApplicationVersion + "." + SCMId
}
