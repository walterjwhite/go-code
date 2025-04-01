package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type data struct {
	value int
}

func main() {
	print("This is a test")
	print([]byte("This is a test"))
	print(&data{value: 1})
}

func print(message interface{}) {
	data, err := json.Marshal(message)
	logging.Panic(err)

	log.Warn().Msgf("Message %s, message: %v\n", data, data)
}
