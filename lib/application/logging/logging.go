package logging

import (
	"fmt"
	"github.com/rs/zerolog/log"

	"runtime/debug"
)

func Error(err error, contextuals ...interface{}) {
	if err == nil {
		return
	}

	if contextuals != nil || len(contextuals) > 0 {
		for i := range contextuals {
			log.Warn().Interface(fmt.Sprintf("contextual: %d", i), contextuals[i]).Msg("Contextual")
		}
	}

	log.Panic().Err(err).Msg("Error")
}

func Warn(err error, message string) {
	if err == nil {
		return
	}

	log.Warn().Msgf("%s - %s", message, err.Error())
	stackTrace := debug.Stack()
	log.Warn().Msgf("Stack trace:\n%s", stackTrace)
}
