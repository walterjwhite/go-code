package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func (e *Executor) MessageDeserialized(deserialized []byte) {
	e.MessageDeserializedWithMetadata(deserialized, &MessageMetadata{})
}

func (e *Executor) MessageDeserializedWithMetadata(deserialized []byte, metadata *MessageMetadata) {
	args, workingDir, err := e.parseMessage(deserialized)
	if err != nil {
		log.Error().Msgf("failed to parse message: %v", err)
		e.publishEvent(metadata, map[string]string{"status": "parse failed", "done": "true", "exitCode": "1", "output": err.Error()})
		return
	}

	if !e.validateAndExecute(args, workingDir, metadata) {
		return
	}
}

func (e *Executor) parseMessage(deserialized []byte) ([]string, string, error) {
	args, workingDir, err := e.handleTarballPayload(deserialized)
	if err == nil {
		return args, workingDir, nil
	}

	log.Debug().Msgf("message is not a tarball, trying to parse as shell command")
	args, err = e.parseShellCommand(deserialized)
	if err != nil {
		return nil, "", err
	}

	return args, "", nil
}

func (e *Executor) validateAndExecute(args []string, workingDir string, metadata *MessageMetadata) bool {
	if len(args) == 0 {
		log.Warn().Msg("no args received")
		e.publishEvent(metadata, map[string]string{"status": "validation failed", "done": "true", "exitCode": "1", "output": "no args received"})
		return false
	}

	if !isValidCommandName(args[0]) {
		log.Warn().Msgf("invalid function name: %s", args[0])
		e.publishEvent(metadata, map[string]string{"status": "validation failed", "done": "true", "exitCode": "2", "output": "invalid function name"})
		return false
	}

	e.publishEvent(metadata, map[string]string{"status": "running cmd"})

	var profile string
	if metadata.Profile != "" {
		profile = metadata.Profile
	}

	status, output := e.executeCommand(args, workingDir, profile, metadata)
	e.publishEvent(metadata, map[string]string{"status": "execution done", "done": "true", "exitCode": string(rune(status)), "output": output})

	return true
}

func (e *Executor) executeCommand(args []string, workingDir string, profile string, metadata *MessageMetadata) (int, string) {
	log.Info().Msgf("running: %s", args)
	ecmd := exec.Command(*cmd, args...)

	if workingDir != "" {
		ecmd.Dir = workingDir
		defer cleanupTempDir(workingDir)
	}

	if profile != "" {
		log.Info().Msgf("using profile %s", profile)
		ecmd.Env = append(os.Environ(), "remote_exec_profile="+profile)
	}

	var b bytes.Buffer
	ecmd.Stdout = &b
	ecmd.Stderr = &b

	if err := ecmd.Start(); err != nil {
		log.Error().Err(err).Msg("Failed to start command")
		return -1, ""
	}

	log.Info().Msgf("Process started: %d", ecmd.Process.Pid)
	e.publishEvent(metadata, map[string]string{"pid": fmt.Sprintf("%d", ecmd.Process.Pid)})

	err := ecmd.Wait()

	output := b.String()
	status := 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		} else {
			status = -1 // System error (e.g., I/O issues)
		}
	}

	return status, output
}

func (e *Executor) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}
