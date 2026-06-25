package main

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (e *Executor) handleTarballPayload(deserialized []byte) ([]string, string, error) {
	tmpDir, err := os.MkdirTemp("", "pubsub-exec-*")
	if err != nil {
		log.Error().Msgf("failed creating temp dir: %v", err)
		return nil, "", err
	}

	if err := extractTarball(tmpDir, deserialized); err != nil {
		log.Error().Msgf("failed extracting tarball to %s: %v", tmpDir, err)
		return nil, "", err
	}

	scriptPath := filepath.Join(tmpDir, "script.sh")
	if err := os.Chmod(scriptPath, 0700); err != nil {
		log.Error().Msgf("failed chmod script %s: %v", scriptPath, err)
		return nil, "", err
	}

	absolutePath, err := filepath.Abs(scriptPath)
	logging.Warn(err, "failed to get absolute path to script")

	return []string{"script_exec", absolutePath}, tmpDir, nil
}

func cleanupTempDir(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Warn().Msgf("failed removing temp dir %s: %v", tmpDir, err)
	}
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

	return nil
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
