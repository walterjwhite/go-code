package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

type ShutdownHandler interface {
	OnShutdown()
	OnContextClosed()
}

var (
	shutdownHooksGroup = sync.WaitGroup{}
)

func Add(ctx context.Context, shutdownHandler ShutdownHandler) {
	shutdownHooksGroup.Add(1)

	osSignalChannel := make(chan os.Signal, 1)

	signal.Notify(osSignalChannel, os.Interrupt)

	go handleShutdown(ctx, osSignalChannel, shutdownHandler)

}

func handleShutdown(ctx context.Context, osSignalChannel chan os.Signal, shutdownHandler ShutdownHandler) {
	select {
	case <-osSignalChannel:
		shutdownHandler.OnShutdown()
	case <-ctx.Done():
		shutdownHandler.OnContextClosed()
	}

	shutdownHooksGroup.Done()
}

/*
func cleanup() {
	shutdownHooksGroup.Wait()
}
*/
