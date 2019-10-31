package remote

import (
"github.com/rs/zerolog/log"
"os"
)

func IsRemote() bool {
	clientName, exists := os.LookupEnv("CLIENTNAME")
	if exists {
		log.Info().Msgf("remote client name: %v", clientName)
	}
	
	return exists
}
