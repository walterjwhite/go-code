package shutdown

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"os"
	"os/signal"
	"sync"
)

type ShutdownHandler interface {
	OnShutdown()
	OnContextClosed()
}

var shutdownHooksGroup = sync.WaitGroup{}
var registerContextCleanupOnce sync.Once

func Add(shutdownHandler ShutdownHandler) {
	shutdownHooksGroup.Add(1)

	osSignalChannel := make(chan os.Signal, 1)
	signal.Notify(osSignalChannel, os.Interrupt)

	go handleShutdown(osSignalChannel, shutdownHandler)

	go registerContextCleanupOnce.Do(cleanup)
}

func handleShutdown(osSignalChannel chan os.Signal, shutdownHandler ShutdownHandler) {
	select {
	case <-osSignalChannel:
		shutdownHandler.OnShutdown()
	case <-application.Context.Done():
		shutdownHandler.OnContextClosed()
	}

	shutdownHooksGroup.Done()
}

func cleanup() {
	shutdownHooksGroup.Wait()
	application.Cancel()
}
