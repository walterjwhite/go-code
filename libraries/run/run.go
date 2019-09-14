package run

import (
	"context"
	"os/exec"
	//"libraries/notify"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func Run(ctx context.Context, applications []string) {
	c := initialize(applications)
	waitForAll(runAll(ctx, c))
}

func initialize(applications []string) Configuration {
	c := Configuration{}

	c.Applications = make([]Application, len(applications))
	for index, application := range applications {
		var a Application
		a.getConf(application)

		c.Applications[index] = a
	}

	return c
}

func runAll(ctx context.Context, c Configuration) []exec.Cmd {
	commands := make([]exec.Cmd, 0)
	for index, application := range c.Applications {
		commands = append(commands, *runApplication(ctx, index, application))
	}

	return commands
}

func waitForAll(commands []exec.Cmd) {
	for _, command := range commands {
		_, err := command.Process.Wait()

		// TODO: push event to channel
		logging.Panic(err)
	}
}
