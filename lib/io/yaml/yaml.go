package yaml

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"os"
)

func Read(configurationFile string, out interface{}) (readErr error) {
	log.Debug().Msgf("reading: %v / %v", configurationFile, out)

	yamlFile, readErr := os.ReadFile(configurationFile)
	if readErr != nil {
		return readErr
	}

	defer func() {
		if r := recover(); r != nil {
			if panicErr, ok := r.(error); ok {
				readErr = panicErr
			} else {
				readErr = fmt.Errorf("yaml.Unmarshal panicked: %v", r)
			}
		}
	}()

	readErr = yaml.Unmarshal(yamlFile, out)
	log.Debug().Msgf("Read:\n%v", out)

	return readErr
}

func Write(in interface{}, outFile string) (writeErr error) {
	var buf []byte

	defer func() {
		if r := recover(); r != nil {
			if panicErr, ok := r.(error); ok {
				writeErr = panicErr
			} else {
				writeErr = fmt.Errorf("yaml.Marshal panicked: %v", r)
			}
		}
	}()

	buf, writeErr = yaml.Marshal(in)
	if writeErr != nil {
		return // Return the error from Marshal or panic
	}

	writeErr = os.WriteFile(outFile, buf, 0644)
	log.Debug().Msgf("Wrote:\n%v", outFile)

	return // Return writeErr
}
