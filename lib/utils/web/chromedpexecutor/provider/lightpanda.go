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
	portAllocation, err := freeport.GetPortAllocation()
	logging.Error(err)

	go runLightpanda(pctx, portAllocation)
	time.Sleep(lightpandaInitDelay)

	return portAllocation.Port
}

func runLightpanda(pctx context.Context, portAllocation *freeport.PortAllocation) {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, lightpandaCmd, "serve", "--host", "127.0.0.1", "--port", strconv.Itoa(portAllocation.Port))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logging.Error(portAllocation.Close())
	logging.Error(cmd.Start())

	log.Info().Msgf("launched lightpanda: 127.0.0.1:%d", portAllocation.Port)

	go shutdownLightpandaGracefully(ctx, cmd)
	go lightpandaExited(cmd, ctx, cancel)
}

func shutdownLightpandaGracefully(ctx context.Context, cmd *exec.Cmd) {
	<-ctx.Done()
	if cmd.Process != nil {
		log.Info().Msg("terminating lightpanda process...")
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Warn().Err(err).Msg("failed to send interrupt to lightpanda, trying kill")
			_ = cmd.Process.Kill()
		}

		_, _ = cmd.Process.Wait()
	}
}

func lightpandaExited(cmd *exec.Cmd, ctx context.Context, cancel context.CancelFunc) {
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == nil {
			log.Error().Err(err).Msg("lightpanda process exited with error")
		}
	} else if ctx.Err() == nil {
		log.Info().Msg("lightpanda process exited successfully")
	}

	cancel()
}
