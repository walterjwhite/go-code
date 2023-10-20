package run

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"

	"errors"

	"bufio"
	"flag"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

var (
	// randomWait                    = flag.Bool("w", false, "introduce random waits between actions")
	detachFromBrowserWhenComplete = flag.Bool("d", true, "detach from browser session when complete")
	sessionFile                   = flag.String("session-file", "", "file to execute")
	chromedpsession               session.ChromeDPSession
)

func Run(ctx context.Context) {
	if len(*sessionFile) == 0 {
		logging.Panic(errors.New("Session File is required"))
	}

	chromedpsession = setup(ctx)

	// if !*randomWait {
	// 	// no need to wait
	// 	chromedpsession.Waiter = nil
	// }

	session.Execute(chromedpsession, ParseActions(read()...)...)
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
	log.Info().Msgf("read: %v", lines)

	return lines
}

func setup(ctx context.Context) session.ChromeDPSession {
	if *detachFromBrowserWhenComplete {
		return remote.New(context.Background())
	}

	return remote.New(ctx)
}
