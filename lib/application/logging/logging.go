package logging

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func Error(err error, contextuals ...any) {
	if err == nil {
		return
	}

	if contextuals != nil || len(contextuals) > 0 {
		for i := range contextuals {
			log.Warn().Interface(fmt.Sprintf("contextual: %d", i), contextuals[i]).Msg("Contextual")
		}
	}

	log.Error().Msgf("error - %s", err.Error())

	if os.Getenv("ENVIRONMENT") == "development" {
		log.Debug().Msgf("Stack trace unavailable in production for security reasons")
	}

	log.Panic().Err(err).Msg("Error")
}

func Warn(err error, message string) {
	if err == nil {
		return
	}

	log.Warn().Msgf("%s - %s", message, err.Error())
}
