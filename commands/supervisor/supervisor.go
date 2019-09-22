package main

import (
	"./irpc/server"
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const refreshInterval = 1 * time.Minute

func refreshAll() error {
	log.Print("refreshing")
	if server.IsEnabled() {
		log.Print("supervision is enabled")
		server.Refresh()
	} else {
		log.Print("supervision is disabled")
	}

	return nil
}

// periodically update the status
func periodic(ctx context.Context, fn func() error, runOnce bool) {
	timer := time.NewTimer(refreshInterval)
	defer timer.Stop()

	// do the initial call
	if runOnce {
		run(fn)
	}

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			run(fn)
			
			// after we're done with this iteration, schedule a new one
			periodic(ctx, refreshAll, false)
		}
	}
}

func run(fn func() error) {
	if err := fn(); err != nil {
		log.Fatalf("Error executing Periodic %v", err)
	}
}

func main() {
	ctx, _ := context.WithCancel(context.Background())
	go periodic(ctx, refreshAll, true)

	rpc.Register(new(server.Server))
	rpc.Register(new(server.InterfaceServer))
	rpc.Register(new(server.DiskServer))
	rpc.Register(new(server.ServiceServer))

	rpc.HandleHTTP()

	l, e := net.Listen(server.Protocol, server.ListenHost+":"+server.Port)
	if e != nil {
		log.Fatalf("Error starting server: %s", e)
	}

	http.Serve(l, nil)
}
