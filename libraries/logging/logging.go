package logging

import (
	"github.com/rs/zerolog/log"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func Warn(err error, isError bool, message ...string) {
	if err != nil {
		if isError {
			log.Error().Msgf("%v", message)
			Panic(err)
		} else {
			log.Warn().Msgf("%v", message)
			log.Warn().Msg(err.Error())
		}
	}
}
