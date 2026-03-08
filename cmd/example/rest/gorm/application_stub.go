package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	Context context.Context
	cancelFunc context.CancelFunc
)

func Configure() {
	Context, cancelFunc = context.WithCancel(context.Background())
}

func OnPanic() {
	if r := recover(); r != nil {
		panic(r)
	}
}

func Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancelFunc()
}
