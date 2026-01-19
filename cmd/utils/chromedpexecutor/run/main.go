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
	}

	sessionDuration, err := time.ParseDuration(*timeString)
	logging.Error(err)

	lines, err := readActions()
	logging.Error(err)

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

	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}
	defer close(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func close(f *os.File) {
	logging.Error(f.Close())
}

func onTabClosed() {
	application.Cancel()
}
