package client

import (
	"encoding/json"
	"encoding/xml"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseXml(t *testing.T) {
	file, err := os.Open("sample.message.xml")
	logging.Panic(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	logging.Panic(err)

	message := &Message{}

	logging.Panic(xml.Unmarshal(data, message))

	log.Info().Msgf("parsed: %v", message)
}

func TestParseJson(t *testing.T) {
	file, err := os.Open("sample.message.json")
	logging.Panic(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	logging.Panic(err)

	message := &Message{}

	logging.Panic(json.Unmarshal(data, message))

	log.Info().Msgf("parsed: %v", message)
}

func TestParseJsonResponse(t *testing.T) {
	file, err := os.Open("sample.response.json")
	logging.Panic(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	logging.Panic(err)

	container := &Container{}

	logging.Panic(json.Unmarshal(data, container))

	log.Info().Msgf("parsed: %v", container.Response.FeedMessageResponse)
	log.Info().Msgf("parsed: %v", container.Response.FeedMessageResponse.Feed)

	for i, message := range container.Response.FeedMessageResponse.Messages.Message {
		log.Info().Msgf("parsed: %v -> %v", i, message)
		log.Info().Msgf("parsed time: %v -> %v/%v/%v %v:%v:%v", i, message.DateTime.Time().Year(), message.DateTime.Time().Month(),
			message.DateTime.Time().Day(), message.DateTime.Time().Hour(), message.DateTime.Time().Minute(), message.DateTime.Time().Second())
	}
}
