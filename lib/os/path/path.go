package path

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
	"path/filepath"
	"strings"
)

var SessionDirectory string

func WithSessionDirectory(sessionDirectory string) {
	SessionDirectory = sessionDirectory
}

func GetFile(label string, extension string, details ...string) *os.File {
	filename := getFilenamePath(getFilename(label, extension, details...))
	directory := getDirectory(filename)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Debug().Msgf("Creating directory: %v", directory)
		logging.Panic(os.MkdirAll(directory, os.ModePerm))

		/*
			if err != nil {
				log.Printf("Error creating directory: %v", filename)
				log.Printf("Error creating directory (dir): %v", filepath.Dir(filename))
				panic(err)
			}
		*/
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)

	return file
}

func getDirectory(filename string) string {
	return filepath.Dir(filename)
}

func getFilenamePath(filename string) string {
	filename, err := homedir.Expand(filename)
	logging.Panic(err)

	return filename
}

func getFilename(label string, extension string, details ...string) string {
	if len(SessionDirectory) == 0 {
		logging.Panic(errors.New("Session Directory was not initialized"))
	}

	filenameWithPrefix := getFilenameWithPrefix(extension, details...)
	return filepath.Join(SessionDirectory, label, filenameWithPrefix)
}

func getFilenameWithPrefix(extension string, details ...string) string {
	if len(details) > 0 {
		return strings.Join(details, ".") + "." + extension
	}

	return extension
}
