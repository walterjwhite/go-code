package logging

import (
	"github.com/rs/zerolog/log"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func Warn(err error, isError bool) {
	if err != nil {
		if isError {
			Panic(err)
		} else {
			log.Warn().Msg(err.Error())
		}
	}
}
