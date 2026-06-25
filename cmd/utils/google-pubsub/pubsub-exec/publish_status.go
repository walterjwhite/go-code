package main

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (e *Executor) publishEvent(metadata *MessageMetadata, data map[string]string) {
	attributes := map[string]string{}
	if metadata.ID != "" {
		attributes["response-to"] = metadata.ID
	}

	attributes["index"] = fmt.Sprintf("%d", metadata.index)
	metadata.index++

	messageBody, err := json.Marshal(data)
	if err != nil {
		log.Warn().Msgf("failed to marshal data: %v", err)
		return
	}

	log.Debug().Msgf("publishing event: %d -> %s", metadata.index, string(messageBody))
	logging.Warn(e.Runtime.Publish(subscriberConf.StatusTopicName, messageBody, attributes), "failed publishing response")
}
