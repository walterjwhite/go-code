package logging

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func logContextuals(contextuals ...any) {
	if len(contextuals) == 0 {
		return
	}
	for i := range contextuals {
		log.Warn().Interface(fmt.Sprintf("contextual: %d", i), contextuals[i]).Msg("Contextual")
	}
}

func isDevEnvironment() bool {
	return os.Getenv("ENVIRONMENT") == "development"
}

func logErrorMessage(err error) {
	log.Error().Msgf("error - %s", err.Error())
}

func logSecurityNote() {
	log.Debug().Msgf("Stack trace unavailable in production for security reasons")
}

func Error(err error, contextuals ...any) {
	if err == nil {
		return
	}

	logContextuals(contextuals...)
	logErrorMessage(err)

	if isDevEnvironment() {
		logSecurityNote()
	}
}

func Warn(err error, message string) {
	if err == nil {
		return
	}

	log.Warn().Msgf("%s - %s", message, err.Error())
}
