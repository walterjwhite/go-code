package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"

	oexec "os/exec"
)

func (e *Executor) MessageDeserialized(deserialized []byte) {
	args, workingDir := e.parseMessage(deserialized)
	if args == nil {
		return
	}

	if !e.validateAndExecute(args, workingDir) {
		return
	}
}

func (e *Executor) parseMessage(deserialized []byte) ([]string, string) {
	var args []string
	if err := json.Unmarshal(deserialized, &args); err == nil {
		return args, ""
	}

	return e.handleTarballPayload(deserialized)
}

func (e *Executor) handleTarballPayload(deserialized []byte) ([]string, string) {
	log.Warn().Msgf("error converting to []string, treating message as tarball payload: %v", json.Unmarshal(deserialized, &[]string{}))

	tmpDir, err := os.MkdirTemp("", "pubsub-exec-*")
	if err != nil {
		log.Error().Msgf("failed creating temp dir: %v", err)
		return nil, ""
	}

	defer cleanupTempDir(tmpDir)

	if err := extractTarball(tmpDir, deserialized); err != nil {
		log.Error().Msgf("failed extracting tarball to %s: %v", tmpDir, err)
		return nil, ""
	}

	scriptPath := filepath.Join(tmpDir, "script.sh")
	if err := os.Chmod(scriptPath, 0700); err != nil {
		log.Error().Msgf("failed chmod script %s: %v", scriptPath, err)
		return nil, ""
	}

	absolutePath, err := filepath.Abs(scriptPath)
	logging.Warn(err, "failed to get absolute path to script")

	return []string{"script_exec", absolutePath}, tmpDir
}

func cleanupTempDir(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Warn().Msgf("failed removing temp dir %s: %v", tmpDir, err)
	}
}

func (e *Executor) validateAndExecute(args []string, workingDir string) bool {
	if len(args) == 0 {
		log.Warn().Msg("no args received")
		return false
	}

	if !isValidCommandName(args[0]) {
		log.Warn().Msgf("invalid function name: %s", args[0])
		return false
	}

	status, output := e.executeCommand(args, workingDir)
	respond(status, output)
	return true
}

func (e *Executor) executeCommand(args []string, workingDir string) (int, string) {
	log.Info().Msgf("running: %s", args)

	ecmd := oexec.Command(*cmd, args...) // #nosec
	if workingDir != "" {
		ecmd.Dir = workingDir
	}

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

	return status, string(output)
}

func (e *Executor) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}

type data struct {
	Status int    `json:"status"`
	Output string `json:"output"`
}

func respond(status int, output string) {
	jsonData, _ := json.Marshal(data{Status: status, Output: output})
	log.Debug().Msgf("response published with status: %v", jsonData)

	logging.Warn(subscriberConf.PubSubConf.Publish(subscriberConf.StatusTopicName, jsonData), "respond")
}

func extractTarball(destination string, payload []byte) error {
	reader := tar.NewReader(bytes.NewReader(payload))

	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := extractTarEntry(destination, header, reader); err != nil {
			return err
		}
	}
}

func extractTarEntry(destination string, header *tar.Header, reader *tar.Reader) error {
	targetPath, err := tarDestinationPath(destination, header.Name)
	if err != nil {
		return err
	}

	switch header.Typeflag {
	case tar.TypeDir:
		return createDirectory(targetPath, header.Mode)
	case tar.TypeReg:
		return createFile(targetPath, header.Mode, reader)
	default:
		log.Warn().Msgf("skipping unsupported tar entry %s with type %d", header.Name, header.Typeflag)
		return nil
	}
}

func createDirectory(path string, mode int64) error {
	return os.MkdirAll(path, os.FileMode(mode))
}

func createFile(path string, mode int64, reader *tar.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return err
	}
	defer func() {
		logging.Warn(file.Close(), "close file")
	}()

	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	return file.Close()
}

func tarDestinationPath(destination, name string) (string, error) {
	cleanName := filepath.Clean(name)
	if cleanName == "." || cleanName == string(filepath.Separator) {
		return "", os.ErrInvalid
	}

	targetPath := filepath.Join(destination, cleanName)
	if isPathTraversal(destination, targetPath) {
		return "", os.ErrPermission
	}

	return targetPath, nil
}

func isPathTraversal(destination, targetPath string) bool {
	relativePath, err := filepath.Rel(destination, targetPath)
	if err != nil {
		return true
	}
	return relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator))
}
