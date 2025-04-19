package run

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider/remote"

	"errors"

	"bufio"
	"flag"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

var (
	detachFromBrowserWhenComplete = flag.Bool("d", true, "detach from browser session when complete")
	sessionFile                   = flag.String("session-file", "", "file to execute")
)

func Run(pctx context.Context) {
	if len(*sessionFile) == 0 {
		logging.Panic(errors.New("session File is required"))
	}

	ctx, cancel := setup(pctx)
	defer cancel()


	action.Execute(ctx, ParseActions(read()...)...)
}

func read() []string {
	log.Info().Msgf("reading: %v", *sessionFile)

	file, err := os.Open(*sessionFile)
	logging.Panic(err)

	defer logging.Panic(file.Close())

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	logging.Panic(scanner.Err())
	log.Info().Msgf("read: %v", lines)

	return lines
}

func setup(ctx context.Context) (context.Context, context.CancelFunc) {
	if *detachFromBrowserWhenComplete {
		return remote.New(context.Background())
	}

	return remote.New(ctx)
}
