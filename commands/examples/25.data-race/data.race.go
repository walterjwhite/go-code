package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

var (
	Context context.Context
	cancel  context.CancelFunc

	osSignalChannel chan os.Signal
)

func main() {
	Context, cancel = context.WithCancel(context.Background())

	//log.Logger = zerolog.New(zerolog.SyncWriter(os.Stdout)).With().Timestamp().Logger()
	//log.Logger = zerolog.New(zerolog.SyncWriter(zerolog.ConsoleWriter{Out: os.Stdout})).With().Timestamp().Logger()
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	osSignalChannel = make(chan os.Signal, 1)
	signal.Notify(osSignalChannel, os.Interrupt)

	time.AfterFunc(5*time.Second, logOnce)

	go handleShutdown()

	//go log1()
	//go log1()

	go log1()

	<-Context.Done()
}

func log1() {
	index := 0
	for {
		time.Sleep(1 * time.Second)
		log.Debug().Msgf("%v", index)

		index++
	}
}

func logOnce() {
	log.Debug().Msg("logOnce")
}

func handleShutdown() {
	log.Info().Msg("handleShutdown")

	defer cancel()

	select {
	case s := <-osSignalChannel:
		log.Info().Msgf("Shutting down: %s", s)
	case <-Context.Done():
		log.Info().Msg("Context is done")
	}
}
