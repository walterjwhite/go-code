package shutdown

import (
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
)

type ShutdownHandler interface {
	OnShutdown(signal os.Signal)
}

type defaultShutdownHandler struct {
}

func Default() context.Context {
	return Register(&defaultShutdownHandler{})
}

func Register(shutdownHandler ShutdownHandler) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	go onShutdown(shutdownHandler, channel, ctx, cancel)

	return ctx
}

func onShutdown(shutdownHandler ShutdownHandler, channel chan os.Signal, ctx context.Context, cancel context.CancelFunc) {
	select {
	case s := <-channel:
		shutdownHandler.OnShutdown(s)
		signal.Stop(channel)
		cancel()

		// TODO: make this configurable, what exit code do we want?
		os.Exit(1)
	case <-ctx.Done():
		log.Printf("Done\n")
	}
}

func (defaultShutdownHandler *defaultShutdownHandler) OnShutdown(signal os.Signal) {
	log.Printf("Shutting down: %s\n", signal)
}
