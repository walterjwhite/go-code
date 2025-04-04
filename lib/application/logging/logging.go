package logging

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func Panic(err error, contextuals ...interface{}) {
	if err != nil {
		if contextuals != nil || len(contextuals) > 0 {
			for i, c := range contextuals {
				log.Warn().Interface(fmt.Sprintf("contextual: %d", i), c).Msg("Contextual")
			}
		}

		log.Panic().Err(err).Msg("Error")
	}
}

func Warn(err error, isError bool, message ...string) {
	if err == nil {
		return
	}

	if isError {
		log.Error().Msgf("%v", message)
		Panic(err)
		return
	}

	log.Warn().Msgf("%v", message)
	log.Warn().Msg(err.Error())
}
