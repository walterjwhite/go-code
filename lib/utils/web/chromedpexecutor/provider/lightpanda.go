package provider

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/freeport"

	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	lightpandaCmd       = "lightpanda"
	lightpandaInitDelay = 250 * time.Millisecond
)

func newLightpandaAllocator(pctx context.Context) (context.Context, context.CancelFunc) {
	port := launchLightpanda(pctx)
	return chromedp.NewRemoteAllocator(pctx, fmt.Sprintf("ws://127.0.0.1:%d", port), chromedp.NoModifyURL)
}

func launchLightpanda(pctx context.Context) int {
	port, err := freeport.GetRandomUnusedPort()
	logging.Panic(err)

	cmd := exec.CommandContext(pctx, lightpandaCmd, "serve", "--host", "127.0.0.1", "--port", strconv.Itoa(port))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	logging.Panic(err)

	log.Info().Msgf("launched lightpanda: 127.0.0.1:%d", port)
	go func() {
		log.Info().Msg("launched lightpanda - go")

		if err := cmd.Wait(); err != nil {
			log.Error().Msgf("Command finished with error: %v\n", err)
		} else {
			log.Info().Msg("Command finished successfully")
		}
	}()

	time.Sleep(lightpandaInitDelay)

	return port
}
