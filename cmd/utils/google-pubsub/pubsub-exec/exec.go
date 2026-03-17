package main

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"

	oexec "os/exec"
)

func (e *Executor) MessageDeserialized(deserialized []byte) {
	e.Args = []string{}

	err := json.Unmarshal(deserialized, &e.Args)
	if err != nil {
		log.Warn().Msgf("error converting to []string, %v | %v", err, deserialized)
		return
	}

	log.Info().Msgf("running: %s", e.Args)

	if len(e.Args) == 0 {
		log.Warn().Msg("no args received")
		return
	}

	if !isValidCommandName(e.Args[0]) {
		log.Warn().Msgf("invalid function name: %s", e.Args[0])
		return
	}

	if len(e.Args) > 1 {
		for i := 1; i < len(e.Args); i++ {
			if !isValidArgument(e.Args[i]) {
				log.Warn().Msgf("invalid argument detected: %s", e.Args[i])
				return
			}
		}
	}

	ecmd := oexec.Command(*cmd, e.Args...) // #nosec
	output, err := ecmd.Output()

	status := 0
	if err != nil {
		if exitError, ok := err.(*oexec.ExitError); ok {
			status = exitError.ExitCode()
			log.Warn().Msgf("Error running: %s (%s) -> %v", *cmd, e.Args, status)
		}
	} else {
		log.Info().Msgf("Successfully ran: %s (%s) -> %v", *cmd, e.Args, status)
	}

	respond(status, string(output))
}

func (e *Executor) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}

type data struct {
	Status int    `json:"status"`
	Output string `json:"output"`
}

func respond(status int, output string) {
	jsonData, _ := json.Marshal(&data{Status: status, Output: output})
	log.Debug().Msgf("response published with status: %v", jsonData)

	logging.Warn(subscriberConf.PubSubConf.Publish(subscriberConf.ResponseTopicName, []byte(jsonData)), "respond")
}
