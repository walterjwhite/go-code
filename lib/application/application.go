package application

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/property"

	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	appFlags = flag.NewFlagSet("application", flag.ExitOnError)
	Context  context.Context
	Cancel   context.CancelFunc
)

func init() {
	Context, Cancel = context.WithCancel(context.Background())
	configureLogging()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan

		Cancel()
	}()
}

func Configure(configurations ...any) {
	isTest := false
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			isTest = true
			break
		}
	}

	if !isTest {
		logging.Error(appFlags.Parse(os.Args[1:]), "appFlags.Parse")
	}
	Load(configurations...)

	configureLogging()

	logIdentifier()
	logStart()
}

func Load(configurations ...any) {
	for i := range configurations {
		if i, ok := configurations[i].(property.PreLoad); ok {
			i.PreLoad()
		}

		property.Load(ApplicationName, configurations[i])

		if i, ok := configurations[i].(property.PostLoad); ok {
			logging.Error(i.PostLoad(Context))
		}
	}
}

func logStart() {
	log.Info().Msg("Application started")
}

func Wait() {
	<-Context.Done()

	log.Info().Msg("Application Context Done")
}

func OnPanic() {
	if r := recover(); r != nil {
		log.Warn().Msgf("panic: %v", r)
		Cancel()

		return
	}
}
