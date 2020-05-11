package main

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func init() {
	application.Configure()
}

func main() {
	d := "V1d6T25aWUtEN1NIV2FFaAo="

	data, err := base64.StdEncoding.DecodeString(d)
	logging.Panic(err)

	log.Info().Msgf("decoded: %v", string(data))
}
