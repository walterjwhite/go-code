package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func main() {
	source := "123456"
	serialized, err := json.Marshal(source)
	logging.Panic(err)

	log.Warn().Msgf("%s -> %s", source, serialized)

	var deserialized string
	err = json.Unmarshal(serialized, &deserialized)
	logging.Panic(err)

	log.Warn().Msgf("%s -> %s", source, deserialized)
}
