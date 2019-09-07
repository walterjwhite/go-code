package path

import (
	"log"
	"os"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"strings"
)

var SessionDirectory string

func With SessionDirectory(sessionDirectory string) {
	SessionDirectory = sessionDirectory
}

func getFile(label string, extension string, details ... string) *os.File {
	filename := getFilenamePath(getFilename(label, extension, details ...))
	directory := getDirectory(filename)
	
	if _, err := os.Stat(directory);os.IsNotExist(err) {
		log.Printf("Creating directory: %v", directory)
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v", filename)
			log.Printf("Error creating directory (dir): %v", filepath.Dir(filename))
			panic(err)
		}
	}
	
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Error creating file %v", err)
	}
	
	return file
}

func getDirectory(filename string) string {
	return filepath.Dir(filename)
}

func getFilenamePath(filename string) string {
	filename, err := homedir.Expand(filename)
	if err != nil {
		panic(err)
	}
	
	return filename
}

func getFilename(label string, extension string, details ... string) string {
	if len(SessionDirectory) == 0 {
		panic("Session Directory was not initialized")
	}
	
	filenameWithPrefix := getFilenameWithPrefix(extension, details ...)
	return filepath.Join(SessionDirectory, label, filenameWithPrefix)
}

func getFilenameWithPrefix(extension string, details ... string) string {
	if len(details) > 0 {
		return strings.Join(details, ".") + "." + extension
	}
	
	return extension
}

