package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"

	oexec "os/exec"
)

func (e *Executor) MessageDeserialized(deserialized []byte) {
	var args []string

	err := json.Unmarshal(deserialized, &args)
	if err != nil {
		log.Warn().Msgf("error converting to []string, treating message as file payload: %v", err)

		tmpFile, tmpErr := os.CreateTemp("", "pubsub-exec-*")
		if tmpErr != nil {
			log.Error().Msgf("failed creating temp file: %v", tmpErr)
			return
		}
		defer func() {
			if closeErr := tmpFile.Close(); closeErr != nil {
				log.Warn().Msgf("failed closing temp file %s: %v", tmpFile.Name(), closeErr)
			}
		}()

		if _, tmpErr = tmpFile.Write(deserialized); tmpErr != nil {
			log.Error().Msgf("failed writing temp file %s: %v", tmpFile.Name(), tmpErr)
			return
		}

		if tmpErr = tmpFile.Chmod(0700); tmpErr != nil {
			log.Error().Msgf("failed chmod temp file %s: %v", tmpFile.Name(), tmpErr)
			return
		}

		name := tmpFile.Name()
		absolutePathToScript, err := filepath.Abs(name)
		logging.Warn(err, "failed to get absolute path to script")
		args = []string{"script_exec", absolutePathToScript}
	}

	log.Info().Msgf("running: %s", args)

	if len(args) == 0 {
		log.Warn().Msg("no args received")
		return
	}

	if !isValidCommandName(args[0]) {
		log.Warn().Msgf("invalid function name: %s", args[0])
		return
	}


	ecmd := oexec.Command(*cmd, args...) // #nosec
	output, err := ecmd.Output()

	status := 0
	if err != nil {
		if exitError, ok := err.(*oexec.ExitError); ok {
			status = exitError.ExitCode()
			log.Warn().Msgf("Error running: %s (%s) -> %v", *cmd, args, status)
		}
	} else {
		log.Info().Msgf("Successfully ran: %s (%s) -> %v", *cmd, args, status)
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

	logging.Warn(subscriberConf.PubSubConf.Publish(subscriberConf.StatusTopicName, []byte(jsonData)), "respond")
}
