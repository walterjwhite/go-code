package plugins

import (
	"github.com/walterjwhite/go-application/libraries/property"
)

type ConfigurationReader interface {
	Load(config interface{}, prefix string)
}

var (
	pluginRegistry map[string]ConfigurationReader
)

func init() {
	readerRegistry = make([]ConfigurationReader, 0)

	readerRegistry = append(readerRegistry, &defaultConfigurationReader{})
	readerRegistry = append(readerRegistry, &envConfigurationReader{})
	readerRegistry = append(readerRegistry, &cliConfigurationReader{})
	readerRegistry = append(readerRegistry, &cliConfigurationReader{})

	readerRegistry = append(readerRegistry, &encryptionReader{})
}

func Initialize(taskName string) {
	foreachfile.Execute(root, initializePlugin)

	// determine yaml filename
	// read yaml file name
	property.Load(config, prefix)
}

func initializePlugin(filePath string) {
	parts := strings.Split(filePath, ".")
	extension := strings.ToLower(parts[len(parts)-1])
}
