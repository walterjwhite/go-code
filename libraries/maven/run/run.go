package run

import (
	"context"
	"os/exec"
	//"libraries/notify"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func Run(ctx context.Context, profile string, debug bool /*, notificationBuilder func(notification notify.Notification) notify.Notifier*/) {
	var c Configuration
	c.getConf(profile)

	commands := make([]exec.Cmd, 0)
	for index, application := range c.Applications {
		commands = append(commands, *runApplication(ctx, index, profile, c, application, debug /*, notificationBuilder*/))
	}

	for _, command := range commands {
		_, err := command.Process.Wait()

		// TODO: push event to channel
		logging.Panic(err)
	}
}
