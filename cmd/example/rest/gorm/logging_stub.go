package main

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func Error(err error) {
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func ErrorWithNil() {
	Error(nil)
}

func ErrorWithActual(err error) {
	Error(err)
}

func NewError(msg string) error {
	return errors.New(msg)
}
