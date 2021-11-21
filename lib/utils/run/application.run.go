package run

import (
	"context"
	"os/exec"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/runner"
)

func (a *Application) Run(ctx context.Context) {
	a.command = exec.CommandContext(ctx, a.Command, a.Arguments...)

	// TODO: make a temp directory
	// TODO: inject files in temp directory
	a.command.Dir = /*a.Name*/ a.session.Path

	notificationChannel := make(chan *string)

	runner.WithEnvironment(a.command, true, a.Environment...)

	a.configureLogWatcher(notificationChannel, a.getLogFile(), a.command)

	log.Debug().Msgf("Running Application: %v (%v)", a.Name, a.Command)
	log.Debug().Msgf("Environment: %v", a.Environment)
	log.Debug().Msgf("Arguments: %v", a.Arguments)
	log.Debug().Msgf("Matcher: %v", a.LogMatcher)

	logging.Panic(a.command.Start())

	go a.monitorChannel(ctx, notificationChannel)
	go a.checkPort(ctx)
}
