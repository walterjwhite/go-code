package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/io/densecode"
	"github.com/walterjwhite/go-code/lib/io/densecode/png"
	"github.com/walterjwhite/go-code/lib/io/densecode/terminal"
)

var (
	configuration = &densecode.Configuration{}

	inputFile  = flag.String("file", "", "source data to be transferred")
	outputPath = flag.String("out", "", "if not specified, images will be shown in the console")
)

func init() {
	application.Configure(configuration)
}

func main() {
	defer application.OnPanic()

	data, err := loadInput(configuration.InputText, configuration.InputFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load input")
	}

	result, err := configuration.Encode(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encode")
	}

	if configuration.WritePNG {
		renderer := png.NewPNGRenderer(configuration.OutputPath)
		if err := renderer.Render(result); err != nil {
			log.Fatal().Err(err).Msg("Failed to render PNG")
		}
	}

	if configuration.Terminal {
		renderer := terminal.NewTerminalRenderer(configuration.AdvanceKey)
		if err := renderer.Render(result); err != nil {
			log.Fatal().Err(err).Msg("Terminal display failed")
		}
	}
}

func loadInput(text, path string) ([]byte, error) {
	if text != "" {
		return []byte(text), nil
	}
	if err := validateFilePath(path); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func validateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal detected: path cannot contain .. components")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("file not accessible: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path points to a directory, not a file")
	}
	return nil
}













