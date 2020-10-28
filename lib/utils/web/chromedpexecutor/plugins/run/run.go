package run

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor"

	"errors"

	"bufio"
	"flag"
	"github.com/walterjwhite/go/lib/application/logging"
	"os"
)

var (
	randomWait                    = flag.Bool("w", false, "introduce random waits between actions")
	detachFromBrowserWhenComplete = flag.Bool("d", true, "detach from browser session when complete")
	sessionFile                   = flag.String("session-file", "", "file to execute")
	chromedpsession               *chromedpexecutor.ChromeDPSession
)

func Run(ctx context.Context) {
	if len(*sessionFile) == 0 {
		logging.Panic(errors.New("Session File is required"))
	}

	chromedpsession = setup(ctx)

	if !*randomWait {
		// no need to wait
		chromedpsession.Waiter = nil
	}

	chromedpsession.Execute(ParseActions(read()...)...)
}

func read() []string {
	log.Info().Msgf("reading: %v", *sessionFile)

	file, err := os.Open(*sessionFile)
	logging.Panic(err)

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	logging.Panic(scanner.Err())

	return lines
}

func setup(ctx context.Context) *chromedpexecutor.ChromeDPSession {
	if *detachFromBrowserWhenComplete {
		return chromedpexecutor.New(context.Background())
	}
	return chromedpexecutor.New(ctx)
}
