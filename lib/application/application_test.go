package application

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"testing"
	"time"
)

func TestConfigure(t *testing.T) {
	Configure()
}

func TestOnPanic(t *testing.T) {
	localCtx, localCancel := context.WithCancel(context.Background())
	Context = localCtx
	Cancel = localCancel

	go func() {
		defer OnPanic()
		panic("test panic")
	}()

	<-Context.Done()
	Context, Cancel = context.WithCancel(context.Background())
}

func TestWait(t *testing.T) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		Cancel()
	}()

	Wait()
}

func TestMain(m *testing.M) {
	originalWriter := log.Logger
	log.Logger = zerolog.New(io.Discard)
	defer func() { log.Logger = originalWriter }()

	os.Exit(m.Run())
}
