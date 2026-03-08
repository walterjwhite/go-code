package property

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/yaml"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

const (
	propertyConfigurationLocation = "~/.config/walterjwhite/go"
)

var (
	getFileFunc = getFile
)

func sanitizeApplicationName(name string) string {
	sanitized := strings.ReplaceAll(name, "..", "")
	sanitized = strings.ReplaceAll(sanitized, "/", "")
	sanitized = strings.ReplaceAll(sanitized, "\\", "")
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	sanitized = reg.ReplaceAllString(sanitized, "")
	return sanitized
}

func LoadFile(applicationName string, config any) error {
	return LoadFileWithPath(config, getFileFunc(applicationName, config))
}

func LoadFileWithPath(config any, filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Warn().Msgf("failed to resolve absolute path: %v", err)
		return err
	}

	cleanPath := filepath.Clean(absPath)

	finfo, err := os.Stat(cleanPath)
	if err != nil {
		log.Warn().Msgf("file does not exist: %v", cleanPath)
		return nil
	}

	if finfo.IsDir() {
		log.Warn().Msgf("file is a directory: %v", cleanPath)
		return nil
	}

	log.Info().Msgf("Reading %v", cleanPath)
	return yaml.Read(cleanPath, config)
}

func getFile(applicationName string, config any) string {
	if len(applicationName) == 0 {
		log.Warn().Msgf("application name is empty: %s", applicationName)
	}

	sanitizedAppName := sanitizeApplicationName(applicationName)
	if sanitizedAppName != applicationName && len(applicationName) > 0 {
		log.Warn().Msgf("application name was sanitized from '%s' to '%s' for security", applicationName, sanitizedAppName)
	}

	path, err := homedir.Expand(propertyConfigurationLocation)
	logging.Error(err, "homedir.Expand")

	return filepath.Join(path, sanitizedAppName, typename.Get(config)+".yaml")
}
