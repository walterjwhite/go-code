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
	logging.Error(err)

	ctx, cancel := context.WithCancel(pctx)

	cmd := exec.CommandContext(ctx, lightpandaCmd, "serve", "--host", "127.0.0.1", "--port", strconv.Itoa(port))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	logging.Error(err)

	log.Info().Msgf("launched lightpanda: 127.0.0.1:%d", port)

	go lightpandaHandleContextEnded(ctx, cmd, cancel)
	go lightpandaExited(cmd, ctx)

	time.Sleep(lightpandaInitDelay)

	return port
}

func lightpandaHandleContextEnded(ctx context.Context, cmd *exec.Cmd, cancel context.CancelFunc) {
	<-ctx.Done()
	if cmd.Process != nil {
		log.Info().Msg("terminating lightpanda process...")
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Warn().Err(err).Msg("failed to send interrupt to lightpanda, trying kill")
			_ = cmd.Process.Kill()
		}

		_, _ = cmd.Process.Wait()
	}

	cancel()
}

func lightpandaExited(cmd *exec.Cmd, ctx context.Context) {
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == nil {
			log.Error().Err(err).Msg("lightpanda process exited with error")
		}
	} else if ctx.Err() == nil {
		log.Info().Msg("lightpanda process exited successfully")
	}
}
