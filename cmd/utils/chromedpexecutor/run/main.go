package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"

	"bufio"
	"errors"
	"flag"
	"path/filepath"
	"strings"

	"os"
	"time"
)

var (
	filename     = flag.String("f", "", "filename to execute")
	timeString   = flag.String("t", "5m", "session timeout")
	providerConf = &provider.Conf{}
)

func main() {
	defer application.OnPanic()

	application.Configure(providerConf)

	if len(*filename) == 0 {
		logging.Error(errors.New("filename is required"))
		return
	}

	sessionDuration, err := time.ParseDuration(*timeString)
	if err != nil {
		logging.Error(err)
		return
	}

	lines, err := readActions()
	if err != nil {
		logging.Error(err)
		return
	}

	ctx, cancel := provider.New(providerConf, application.Context)
	defer cancel()



	action.OnTabClosed(ctx, onTabClosed)

	logging.Error(action.Execute(ctx, run.ParseActions(lines...)...))

	select {
	case <-time.After(sessionDuration):
		log.Info().Msgf("session duration limit reached: %v", sessionDuration)
	case <-application.Context.Done():
		log.Info().Msg("application context done")
	}
}

func readActions() ([]string, error) {
	var lines []string

	cleanPath := filepath.Clean(*filename)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, errors.New("invalid file path: unable to resolve absolute path")
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("unable to determine working directory")
	}

	relPath, err := filepath.Rel(wd, absPath)
	if err != nil {
		return nil, errors.New("invalid file path: unable to resolve relative path")
	}

	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return nil, errors.New("path traversal detected: file must be within the working directory")
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, err
	}
	defer closeResource(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func closeResource(f *os.File) {
	logging.Warn(f.Close(), "failed to close file")
}

func onTabClosed() {
	application.Cancel()
}
