package identifier

import (
	"log"
)

// *MUST* be injected at compile time
var applicationVersion string
var buildDate string

func Log() {
	if isConfigured() {
		log.Printf("Application Version: %v\n", applicationVersion)
		log.Printf("Built on: %v\n", buildDate)
	}
}

func isConfigured() bool {
	return len(applicationVersion) > 0 && len(buildDate) > 0
}
