package sudo

import (
	"context"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/runner"
	"os/user"
)

func Run(ctx context.Context, command string, arguments ...string) (int, error) {
	if !isRoot() {
		return runner.Run(ctx, "sudo", append(arguments[:0], append([]string{command}, arguments[1:]...)...)...)
	}

	return runner.Run(ctx, command, arguments...)
}

func isRoot() bool {
	return getUser() == "root"
}

func getUser() string {
	user, err := user.Current()
	logging.Panic(err)

	return user.Username
}
